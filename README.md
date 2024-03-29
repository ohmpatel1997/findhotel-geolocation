# FindHotel-Geolocation

This project consistes of two cmponents:

1) An API which provides the details of an IP address such as city, country, state, its coordinates etc.
2) An import service which imports the the Geolocation data into the database given a CSV file.

<h1>Installation</h1>

<h2>Prerequisite:</h2>

1) The docker is installed locally.
2) The latest version of Go. You can get it from https://golang.org/doc/install
3) The `make` tool is installed.

<h2> Steps to start the API </h2>
 

1) Start the docker locally
3) Run the command in terminal: `cd && cd findhotel-geolocation/cmd/docker/client-api`
4) Run the command: `make all`

Now you can call the API at `GET http://localhost:3000/v1/ip-info?ip=<ip_address>`

Before calling the API you need to import the data in your database. Please follow below instruction to import the data. Make sure you have already started API by now using above steps.

<h2> Steps to import the data </h2>
  
1) Prepare a csv file named `data_dump.csv` in following format and put it in the `findhotel-geolocation/cmd/docker/import` directory of project.
    
    <h3>CSV Format:</h3>
            
            ip_address,country_code,country,city,latitude,longitude,mystery_value
            200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346

2) Run the command in terminal: `cd && cd findhotel-geolocation/migration`
3) Run the command: `make migrate-up`
3) Run the command: `cd && cd findhotel-geolocation/cmd/docker/import`
4) Run the command: `make all`

Note: Please wait until the import is complete. Depending on your csv file size, it might take time.

After the import is done, you can call the above API, it will work.



------------------------------------------------------------------------------------------------------------------------
You can also access already built API deployed on heroku with already imported data on: 

`GET https://geo-location-assignment.herokuapp.com/v1/ip-info?ip=<ip_address>`

The API will take an IP address as query parameters and give back data for that ip address.
           

