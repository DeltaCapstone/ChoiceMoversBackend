package main

import (
	//"errors"
	"errors"
	"fmt"
	"net/http"

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
func itemsToSizes(estRequest models.UnownedEstimateRequest) ([]int, error) {
	itemMap := map[string]int{
		"table": 1,
		"sm":    1,
		"md":    2,
		"lg":    3,
	}
	sizes := []int{0, 0, 0}
	var size int

	for _, room := range estRequest.Rooms {
		for item, quantity := range room {
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

func calculateItemLoad(estRequest models.UnownedEstimateRequest) (int, error) {
	return 0, nil
}

// Represents the amount of hours it takes to pack a box
// 0.25 = 4 boxes packed per hour
var boxMultiplier float64 = 0.25

func packHours(estRequest models.UnownedEstimateRequest, boxes int) (float64, error) {
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

func loadHours(estRequest models.UnownedEstimateRequest, itemLoad int) (float64, error) {
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
func estimateHours(estRequest models.UnownedEstimateRequest, boxes int, itemLoad int) (float64, error) {
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

func estimateWorkers(estRequest models.UnownedEstimateRequest) (int, error) {
	return 2, nil
}

// Calculate the total cost of a estimate
func estimateCost(estRequest models.UnownedEstimateRequest, hours float64, workers int) (float64, error) {

	return 0, nil
}

// Creates an estimate object from an Unowned Estimate. Used by both owned and unowned estimate creation.
// Calculates labor hours, milage, cost, etc. for an estimate.
func calculateEstimate(req models.UnownedEstimateRequest, c echo.Context) (models.Estimate, error) {
	var estimate models.Estimate
	itemLoad, err := calculateItemLoad(req)
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

	sizes, err := itemsToSizes(req)
	if err != nil {
		return estimate, err
	}

	// Calculate Labor Hours
	hours, err := estimateHours(req, 0, itemLoad)
	if err != nil {
		return estimate, err
	}

	hours_interval := pgtype.Interval{}

	workers, err := estimateWorkers(req)
	if err != nil {
		return estimate, err
	}

	// Calculate the cost of the job
	cost, err := estimateCost(req, hours, workers)
	if err != nil {
		return estimate, err
	}

	estimate = models.Estimate{
		LoadAddrID:   loadAddrID,
		UnloadAddrID: unloadAddrID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,

		Rooms:      req.Rooms,
		Special:    req.Special,
		Small:      sizes[0],
		Medium:     sizes[1],
		Large:      sizes[2],
		Boxes:      0,
		ItemLoad:   itemLoad,
		FlightMult: float64(req.Flights),

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
		EstimateRate:     0.0,
		EstimateCost:     cost,
	}

	return estimate, nil
}

// POST Route to create an Estimate with an account
func createEstimate(c echo.Context) error {
	var req models.CreateEstimateRequest
	// attempt at binding incoming json to a jobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	args, err := calculateEstimate(models.UnownedEstimateRequest{
		LoadAddr:   req.LoadAddr,
		UnloadAddr: req.UnloadAddr,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,

		Rooms:   req.Rooms,
		Special: req.Special,
		Boxes:   req.Boxes,
		Flights: req.Flights,

		Pack:   req.Pack,
		Unpack: req.Unpack,
		Load:   req.Load,
		Unload: req.Unload,

		Clean: req.Clean,

		NeedTruck: req.NeedTruck,
	}, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error"})
	}

	args.CustomerUsername = req.Customer.UserName

	est_id, err := DB.PgInstance.CreateEstimate(c.Request().Context(), args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return c.JSON(http.StatusConflict, fmt.Sprintf("Not Null violation: %v", err))
			}
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create estimate: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"estimate id": est_id})
}

// POST route for unauthenticated estimate requests
func createUnownedEstimate(c echo.Context) error {
	var req models.UnownedEstimateRequest
	// attempt at binding incoming json to an Unowned Estimate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	result, err := calculateEstimate(req, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"cost": result.EstimateCost})
}
