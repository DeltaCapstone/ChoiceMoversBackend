--inserting test users -- NO PASSWORDS FOR NOW

-- Inserting test customers with random names
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_john_d', 'John', 'Doe', 'john.doe@test.com', '330-111-1111', 'abcde');
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_jane_s', 'Jane', 'Smith', 'jane.smith@test.com', '330-111-1112', 'fghij');
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_emily_b', 'Emily', 'Brown', 'emily.brown@test.com', '330-111-1113', 'klmno');
INSERT INTO customers (username, first_name, last_name, email, phone_primary, password_hash) VALUES ('user_mike_w', 'Mike', 'Wilson', 'mike.wilson@test.com', '330-111-1114', 'pqrst');

-- Inserting test employees with random names
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_alex_j', 'Alex', 'Johnson', 'Part-time', 'alex.johnson@test.com', '330-111-2222', 'uvwxy');
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_linda_k', 'Linda', 'King', 'Full-time', 'linda.king@test.com', '330-111-2223', 'zabcd');
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('emp_david_l', 'David', 'Lee', 'Full-time', 'david.lee@test.com', '330-111-2224', 'efghi');
INSERT INTO employees (username, first_name, last_name, employee_type, email, phone_primary, password_hash) VALUES ('mgr_sarah_m', 'Sarah', 'Miller', 'Manager', 'sarah.miller@test.com', '330-111-3333', 'jklmn');
