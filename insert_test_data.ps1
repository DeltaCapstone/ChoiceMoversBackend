
docker cp ./sql_scripts/test_data.sql db:/home/test_data.sql

docker exec -it db psql -d choicemovers -U postgres -a -f /home/test_data.sql