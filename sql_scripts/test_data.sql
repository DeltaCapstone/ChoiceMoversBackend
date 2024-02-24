--inserting test users -- NO PASSWORDS FOR NOW
INSERT INTO customers (username,email,phone_primary) VALUES ('customer1','customer1@test.com','330-111-1111');
INSERT INTO customers (username,email,phone_primary) VALUES ('customer2','customer2@test.com','330-111-1112');
INSERT INTO customers (username,email,phone_primary) VALUES ('customer3','customer3@test.com','330-111-1113');
INSERT INTO customers (username,email,phone_primary) VALUES ('customer4','customer4@test.com','330-111-1114');

INSERT INTO employees (username,employee_type,email,phone_primary) VALUES ('employee1','Part-time','employee1@test.com','330-111-2222');
INSERT INTO employees (username,employee_type,email,phone_primary) VALUES ('employee2','Full-time','employee2@test.com','330-111-2223');
INSERT INTO employees (username,employee_type,email,phone_primary) VALUES ('employee3','Full-time','employee3@test.com','330-111-2224');

INSERT INTO employees (username,employee_type,email,phone_primary) VALUES ('manager1','Manager','manager1@test.com','330-111-3333'); 