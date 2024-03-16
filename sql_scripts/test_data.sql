
-- Inserting test customers with random names
INSERT INTO customers 
    (username, first_name, last_name, email, phone_primary, password_hash) 
VALUES 
    ('user_john_d', 'John', 'Doe', 'john.doe@test.com', '330-111-1111', '$2a$10$rpcIFTYGPDVrS4GpQJzWpOuDay540ZWrsWjvzm13YY9/OqTxjkvcS'), --pw: abcde
    ('user_jane_s', 'Jane', 'Smith', 'jane.smith@test.com', '330-111-1112', '$2a$10$4QHNTDYIvHPJPa/uV7rnxe1gGuSg9P4sBYiIctHn6Y79BsmkPaKvq'), --pw: fghij
    ('user_emily_b', 'Emily', 'Brown', 'emily.brown@test.com', '330-111-1113', '$2a$10$ncOlS1qa5ZNWa9OC4k5h7.NujWD1wTcWS5mV.gjy8gSt2EzWP1lGi'), --pw: klmno
    ('user_mike_w', 'Mike', 'Wilson', 'mike.wilson@test.com', '330-111-1114', '$2a$10$KT6/ddMC3o3wMvPYuAUQSu1kBUeKMUGIOwsUKfOMszHlW4Y5e3Vwu'); --pw: pqrst

-- Inserting test employees with random names
INSERT INTO employees 
    (username, first_name, last_name, employee_type, employee_priority, email, phone_primary, password_hash)
VALUES 
    ('alex.j', 'Alex', 'Johnson', 'Part-time',3 ,'alex.johnson@test.com', '330-111-2222', '$2y$10$eACI6jHv5kopEk92l6KVqO63LvyVYxCvWjXta5Cq9AGCeoy2a1Vvq'), --pw: password
    ('linda.k', 'Linda', 'King', 'Full-time', 2,'linda.king@test.com', '330-111-2223', '$2y$10$eACI6jHv5kopEk92l6KVqO63LvyVYxCvWjXta5Cq9AGCeoy2a1Vvq'), --pw: password
    ('david.l', 'David', 'Lee', 'Full-time', 1,'david.lee@test.com', '330-111-2224', '$2y$10$eACI6jHv5kopEk92l6KVqO63LvyVYxCvWjXta5Cq9AGCeoy2a1Vvq'), --pw: password
    ('sarah.m', 'Sarah', 'Miller', 'Manager', 1,'sarah.miller@test.com', '330-111-3333', '$2y$10$eACI6jHv5kopEk92l6KVqO63LvyVYxCvWjXta5Cq9AGCeoy2a1Vvq'); --pw: password


INSERT INTO public.addresses 
    (street, city, state, zip, res_type, flights, apt_num,square_feet)
VALUES 
    ('123 Main St', 'Cleveland', 'OH', '44101', 'House', 2, 'Apt 101',1500),
    ('456 Oak St', 'Akron', 'OH', '44302', 'Apartment', 1, 'Apt 202',1500),
    ('789 Elm St', 'Cuyahoga Falls', 'OH', '44221', 'House', 0, 'Apt 303',1500),
    ('101 Pine St', 'Cleveland Heights', 'OH', '44118', 'House', 3, 'Apt 404',1500),
    ('202 Maple St', 'Euclid', 'OH', '44123', 'Apartment', 0, 'Apt 505',1500),
    ('303 Cedar St', 'Lakewood', 'OH', '44107', 'House', 1, 'Apt 606',1500),
    ('404 Birch St', 'Shaker Heights', 'OH', '44120', 'House', 4, 'Apt 707',1500),
    ('505 Spruce St', 'Parma', 'OH', '44129', 'Apartment', 2, 'Apt 808',1500);


-- Inserting 4 estimates
INSERT INTO public.estimates 
    (customer_username, load_addr_id, unload_addr_id, start_time, end_time, rooms, special, small_items, medium_items, large_items, boxes, item_load, flight_mult, pack, unpack, load, unload, clean, need_truck, number_workers, dist_to_job, dist_move, estimated_man_hours, estimated_rate, estimated_cost) 
VALUES 
    ('user_john_d', 1, 2, '2024-04-10 08:00:00', '2024-04-10 16:00:00', '{"bedroom": 2, "living_room": 1}', '{}', 5, 3, 1, 10, 100, 1.0, true, false, true, false, true, true, 3, 10, 20, '4 hours', 25.00, 100.00),
    ('user_jane_s', 3, 4, '2024-04-15 10:00:00', '2024-04-15 18:00:00', '{"bedroom": 3, "kitchen": 1}', '{}', 7, 4, 2, 15, 120, 1.2, true, true, true, false, false, false, 2, 15, 25, '5 hours', 30.00, 150.00),
    ('user_emily_b', 5, 6, '2024-04-20 09:00:00', '2024-04-20 17:00:00', '{"bedroom": 4, "bathroom": 2}', '{"fragile_items": ["vases", "glasses"]}', 8, 5, 3, 20, 130, 1.5, false, true, false, true, true, true, 4, 20, 30, '6 hours', 35.00, 180.00),
    ('user_mike_w', 7, 8, '2024-04-25 11:00:00', '2024-04-25 19:00:00', '{"living_room": 1, "dining_room": 1}', '{"antiques": ["painting", "sculpture"]}', 6, 4, 2, 12, 110, 1.3, true, false, false, true, false, true, 3, 18, 22, '4.5 hours', 28.00, 130.00);

-- Inserting 4 jobs
INSERT INTO public.jobs 
    (estimate_id, man_hours, rate, cost, finalized, actual_man_hours, final_cost, ammount_payed, notes) 
VALUES 
    (1, '8 hours', 20.00, 160.00, true, '8 hours', 160.00, 150.00, 'Job completed successfully.'),
    (2, '8 hours', 25.00, 200.00, true, '8 hours', 200.00, 180.00, 'Additional items requested by the customer.'),
    (3, '8 hours', 30.00, 240.00, true, '8 hours', 240.00, 220.00, 'Customer provided additional instructions.'),
    (4, '8 hours', 35.00, 280.00, true, '8 hours', 280.00, 250.00, 'Job completed on time and within budget.');
