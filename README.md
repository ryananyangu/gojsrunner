limit size of file of service
limit function inputs and out puts to specific structures. request build --> response build
passdown statics

cache entire request response templates for fater processing
use grpc for faster inter service communication

--> wrappercore --> invoked script (request, statics,headers) --> string and headers object
--> response --> response payload -(response) --> structured response


Recieve request
-> Get the configs (constants , headers ,payload) based on the service id[Done]
-> Get the client credential if available from the db based on client ID and serviceCode[Done]
-> Pass the configs and payload to wrapper script for building including client credentials[Done]
-> Get the response from the script[Done]
-> send post request to the TP service[Done]
-> Get response and convert it to json if not already in json format
-> Pipe to post send request script[Done]
-> Get the response back and propagate it to the caller[Done]
------------------------------------------------------------------------------------------
-> Pass initial payment params back to the response builder for easier mapping of final feedback (Cancelled)
-> Deal with caching (In progress)
-> Deal with database inserts (Done)
-> XML convertion back to json (Inprogress)
-> Service Type information on service config to state sync or async (Done)
-> Callback endpoint on dxl for async (In Progress)
-> Only if sync we post into the callbackQ else we only post in the database with async function
-> Swagger setup on the dxl layer
-> Add metabase for UI
-> Unit tests
-> Load tests
-> Code refactor
-> Error handling and ACK management on the TP Service.
