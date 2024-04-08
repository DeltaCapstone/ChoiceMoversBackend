package main

import (
	//"errors"
	"errors"
	"fmt"
	"net/http"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	//"github.com/jackc/pgerrcode"
	//"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

//TODO: Redo error handling to get rid of of al lthe sprintf's

func listJobs(c echo.Context) error {
	startDate := c.QueryParam("start")
	endDate := c.QueryParam("end")
	var status string
	if c.Get("role").(string) == "Manager" {
		status = "all"
	} else {
		status = "finalized"
	}

	jobs, err := DB.PgInstance.GetJobsByStatusAndRange(c.Request().Context(), status, startDate, endDate)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if jobs == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No jobs found with status: %v", status))
	}
	return c.JSON(http.StatusOK, jobs)
}

// Creates an address and inserts it into the database.
// Returns the address ID
func createAddress(address models.Address, c echo.Context) (int, error) {
	addr_id, err := DB.PgInstance.CreateAddress(c.Request().Context(), address)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return 0, err
			}
		}
		return 0, err
	}
	return addr_id, nil
}

// Maps all items in the json to their corresponding sizes
func itemsToSizes(estRequest models.EstimateRequest) ([]int, error) {
	itemMap := map[string]int{
		"table":     0,
		"poolTable": 2,
		"couch":     2,
		"loofa":     1,
		"lamp":      0,
		"smBox":     0,
		"mdBox":     1,
		"lgBox":     2,
	}

	// sm, md, lg
	sizes := []int{0, 0, 0}
	var size int

	for _, room := range estRequest.Rooms {
		for item, quantity := range room.Items {
			size = itemMap[item]
			sizes[size] += quantity
		}
	}

	for item, quantity := range estRequest.Special {
		size = itemMap[item]
		sizes[size] += quantity
	}

	for item, quantity := range estRequest.Boxes {
		size = itemMap[item]
		sizes[size] += quantity
	}

	return sizes, nil
}

func calculateItemLoad(sizes []int) (float64, error) {
	sm_mult := 0.02
	md_mult := 0.05
	lg_mult := 0.1

	load := (sm_mult*float64(sizes[0]) +
		md_mult*float64(sizes[1]) +
		lg_mult*float64(sizes[2]))

	return load, nil
}

// Represents the amount of hours it takes to pack a box
// 0.25 = 4 boxes packed per hour
var boxMultiplier float64 = 0.25

func packHours(estRequest models.EstimateRequest, boxes int) (float64, error) {
	// Converts Pack and Unpack bools into ints and uses them as the multiplier for the hours
	// If neither are true, no hours will be added for packing
	packMult := 0.0
	if estRequest.Pack {
		packMult += 1
	}
	if estRequest.Unpack {
		packMult += 1
	}

	return (boxMultiplier * float64(boxes) * packMult), nil
}

func loadHours(estRequest models.EstimateRequest, itemLoad float64) (float64, error) {
	// Converts Load and Unload bools into ints and uses them as the multiplier for the hours
	// If neither are true, no hours will be added for loading
	loadMult := 0.0
	if estRequest.Load {
		loadMult += 1
	}
	if estRequest.Unload {
		loadMult += 1
	}

	return float64(itemLoad) * loadMult, nil
}

// Calculates how many hours a estimate should take
func estimateHours(estRequest models.EstimateRequest, boxes int, itemLoad float64) (float64, error) {
	pack, err := packHours(estRequest, boxes)
	if err != nil {
		return 0, err
	}

	load, err := loadHours(estRequest, itemLoad)
	if err != nil {
		return 0, err
	}

	return pack + load, nil
}

func estimateWorkers(estRequest models.EstimateRequest) (int, error) {
	// Maps special item names to their number of needed workers
	specials := map[string]int{
		"poolTable": 3,
		"piano":     4,
	}

	numWorkers := 2

	for item, quantity := range estRequest.Special {
		if quantity < 1 {
			continue
		}

		val, ok := specials[item]
		if ok {
			if val > numWorkers {
				numWorkers = val
			}
		}
	}

	return numWorkers, nil
}

// Calculates the cost of an hour
func estimateRate(estRequest models.EstimateRequest, workers int) (float64, error) {
	distRate := 10 * (float64(estRequest.DistToJob-15) / 15)
	manRate := workers * 40
	rate := distRate + float64(manRate)
	return rate, nil
}

// Calculate the total cost of a estimate
func estimateCost(estRequest models.EstimateRequest, hours float64, workers int, rate float64) (float64, error) {
	jobHours := hours / float64(workers)
	// Assumes that DistMove is in minutes or miles (60 mph)
	totalHours := jobHours + (float64(estRequest.DistMove) / 60)
	return totalHours * rate, nil
}

// Creates an estimate object from an Unowned Estimate. Used by both owned and unowned estimate creation.
// Calculates labor hours, milage, cost, etc. for an estimate.
func calculateEstimate(req models.EstimateRequest, c echo.Context) (models.Estimate, error) {
	var estimate models.Estimate

	sizes, err := itemsToSizes(req)
	if err != nil {
		return estimate, err
	}

	itemLoad, err := calculateItemLoad(sizes)
	if err != nil {
		return estimate, err
	}

	var loadAddrID, unloadAddrID int
	if req.LoadAddr != nil {
		loadAddrID, err = createAddress(*req.LoadAddr, c)
		if err != nil {
			return estimate, err
		}
	} else {
		loadAddrID = -1
	}

	if req.UnloadAddr != nil {
		unloadAddrID, err = createAddress(*req.UnloadAddr, c)
		if err != nil {
			return estimate, err
		}
	} else {
		loadAddrID = -1
	}

	boxes := req.Boxes["smBox"] + 2*req.Boxes["mdBox"] + 4*req.Boxes["lgBox"]
	fmt.Println(req.Boxes)

	// Calculate Labor Hours
	hours, err := estimateHours(req, boxes, itemLoad)
	if err != nil {
		return estimate, err
	}

	hours_interval := pgtype.Interval{
		Microseconds: int64(hours * 3600000000),
		Valid:        true,
	}

	workers, err := estimateWorkers(req)
	if err != nil {
		return estimate, err
	}

	rate, err := estimateRate(req, workers)
	if err != nil {
		return estimate, err
	}

	// Calculate the cost of the job
	cost, err := estimateCost(req, hours, workers, rate)
	if err != nil {
		return estimate, err
	}

	startTime, err := time.Parse("02.01.2006 03:04 PM", req.StartTime)
	if err != nil {
		return estimate, err
	}

	timeStamp := pgtype.Timestamp{
		Time:  startTime,
		Valid: true,
	}

	estimate = models.Estimate{
		LoadAddrID:   loadAddrID,
		UnloadAddrID: unloadAddrID,
		StartTime:    timeStamp,

		Rooms:      req.Rooms,
		Special:    req.Special,
		Small:      sizes[0],
		Medium:     sizes[1],
		Large:      sizes[2],
		Boxes:      boxes,
		ItemLoad:   itemLoad,
		FlightMult: float64((*req.UnloadAddr).Flights),

		Pack:   req.Pack,
		Unpack: req.Unpack,
		Load:   req.Load,
		Unload: req.Unload,

		Clean: req.Clean,

		NeedTruck:     req.NeedTruck,
		NumberWorkers: workers,
		DistToJob:     req.DistToJob,
		DistMove:      req.DistMove,

		EstimateManHours: hours_interval,
		EstimateRate:     rate,
		EstimateCost:     cost,
	}

	return estimate, nil
}

// POST Route to create an Estimate with an account
func createEstimate(c echo.Context) error {
	var req models.EstimateRequest
	// attempt at binding incoming json to a jobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	args, err := calculateEstimate(req, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	args.CustomerUsername = c.Get("username").(string)

	_, err = DB.PgInstance.CreateEstimate(c.Request().Context(), args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return c.JSON(http.StatusConflict, fmt.Sprintf("Not Null violation: %v ----- Data: %v", err, args))
			}
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create estimate: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"result": args})
}

// POST route for unauthenticated estimate requests
func createUnownedEstimate(c echo.Context) error {
	var req models.EstimateRequest
	// attempt at binding incoming json to an Unowned Estimate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	result, err := calculateEstimate(req, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"result": result})
}
