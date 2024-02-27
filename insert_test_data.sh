#!/bin/bash
docker cp ./sql_scripts/test_data.sql db:/home/test_data.sql
#docker cp ./sql_scripts/create_tables.sql db:/home/create_tables.sql

#docker exec -it db psql -d choicemovers -U postgres -a -f /home/create_tables.sql
docker exec -it db psql -d choicemovers -U postgres -a -f /home/test_data.sql
