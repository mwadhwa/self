# self

## This project is based on rest APIs. Following are the steps to run the project

1. Start the server from main.go. Server will be running on "localhost:8080"
2. This project supports following APIs:
   ### Subscribe - subscribe to a new address
   ```
   curl --location --request POST 'localhost:8080/transactions/subscribe' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "address":"0x6b75d8af000000e20b7a7ddf000ba900b4009a80"
   }'
    ```

    ### Get Transactions of a given address:
    ```
    curl --location --request GET 'localhost:8080/transactions/addresses/0x6b75d8af000000e20b7a7ddf000ba900b4009a80'
    ```

    ### Get Current block:
   ```
    curl --location --request GET 'localhost:8080/block'
    ```
   

Resource folder has screenshot of the response for reference.


## Some considerations:
1. In Memory store is safe to use concurrently. However, It assumes that given address will not be overridden by multiple threads.
2. This is very basic implementations and can be extended to scale and modular but due to time constraint kept it like this