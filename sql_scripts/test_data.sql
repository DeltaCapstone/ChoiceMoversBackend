
-- Inserting test customers with random names
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_john_d', 'John', 'Doe', 'john.doe@test.com', '330-111-1111', '$2a$10$rpcIFTYGPDVrS4GpQJzWpOuDay540ZWrsWjvzm13YY9/OqTxjkvcS'); --pw: abcde
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_jane_s', 'Jane', 'Smith', 'jane.smith@test.com', '330-111-1112', '$2a$10$4QHNTDYIvHPJPa/uV7rnxe1gGuSg9P4sBYiIctHn6Y79BsmkPaKvq'); --pw: fghij
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_emily_b', 'Emily', 'Brown', 'emily.brown@test.com', '330-111-1113', '$2a$10$ncOlS1qa5ZNWa9OC4k5h7.NujWD1wTcWS5mV.gjy8gSt2EzWP1lGi'); --pw: klmno
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_mike_w', 'Mike', 'Wilson', 'mike.wilson@test.com', '330-111-1114', '$2a$10$KT6/ddMC3o3wMvPYuAUQSu1kBUeKMUGIOwsUKfOMszHlW4Y5e3Vwu'); --pw: pqrst

-- Inserting test employees with random names
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_alex_j', 'Alex', 'Johnson', 'Part-time', 'alex.johnson@test.com', '330-111-2222', '$2a$10$5CD1oYBcvimXPQOPH035buUP/e4ZfW4qJiabNArF9qobKBfoJuB0m'); --pw: uvwxy
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_linda_k', 'Linda', 'King', 'Full-time', 'linda.king@test.com', '330-111-2223', '$2a$10$t2nKx8UirO9SXUySR79S5.iTjAGaylkTYIz7Fn7IUzj43nTzgPrw2'); --pw: zabcd
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_david_l', 'David', 'Lee', 'Full-time', 'david.lee@test.com', '330-111-2224', '$2a$10$JzeCi35FzkECqJYPvEO4QuGx0hRW5W6kykEFjo1/wsxPWKGs1UHpm'); --pw: efghi
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('mgr_sarah_m', 'Sarah', 'Miller', 'Manager', 'sarah.miller@test.com', '330-111-3333', '$2a$10$mGZwtnxS8SiwGSS1Q2olP.dht8WrCb1o5NRuuANYe4D5LH5hOtmAe'); --pw: jklmn


INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('123 Main St', 'Cleveland', 'OH', '44101', 'House', 2, 'Apt 101');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('456 Oak St', 'Akron', 'OH', '44302', 'Apartment', 1, 'Apt 202');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('789 Elm St', 'Cuyahoga Falls', 'OH', '44221', 'House', 0, 'Apt 303');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('101 Pine St', 'Cleveland Heights', 'OH', '44118', 'House', 3, 'Apt 404');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('202 Maple St', 'Euclid', 'OH', '44123', 'Apartment', 0, 'Apt 505');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('303 Cedar St', 'Lakewood', 'OH', '44107', 'House', 1, 'Apt 606');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('404 Birch St', 'Shaker Heights', 'OH', '44120', 'House', 4, 'Apt 707');

INSERT INTO public.addresses (street, city, state, zip, res_type, flights, apt_num)
VALUES ('505 Spruce St', 'Parma', 'OH', '44129', 'Apartment', 2, 'Apt 808');


INSERT INTO public.jobs (customer_id, load_addr, unload_addr, start_time, hours_labor, finalized, rooms, pack, unpack, load, unload, clean, milage, cost)
VALUES (1, 1, 2, '2024-04-25 10:00:00', '4 hours', false, '{}'::jsonb, true, false, true, true, false, 10, 150.00);

INSERT INTO public.jobs (customer_id, load_addr, unload_addr, start_time, hours_labor, finalized, rooms, pack, unpack, load, unload, clean, milage, cost)
VALUES (2, 3, 4, '2024-04-26 12:30:00', '5 hours', false, '{}'::jsonb, true, true,true, true, true, 15, 200.00);

INSERT INTO public.jobs (customer_id, load_addr, unload_addr, start_time, hours_labor, finalized, rooms, pack, unpack, load, unload, clean, milage, cost)
VALUES (3, 5, 6, '2024-04-27 08:45:00', '3 hours', false, '{}'::jsonb, false, false, true, true, false, 8, 120.00);

INSERT INTO public.jobs (customer_id, load_addr, unload_addr, start_time, hours_labor, finalized, rooms, pack, unpack, load, unload, clean, milage, cost)
VALUES (4, 7, 8, '2024-04-28 14:15:00', '6 hours', false, '{}'::jsonb, true, false, true, false, false, 20, 250.00);
