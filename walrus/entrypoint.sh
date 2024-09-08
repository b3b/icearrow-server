#!/bin/bash

set -xe

function init_sui() {

    # create new enviroment
    sui client --yes active-address

    active_env=$(sui client active-env)
    if [[ "$active_env" == "testnet" ]]; then
        # Get the active address
        active_address=$(sui client --yes active-address)

        # Check if the address starts with '0x'
        if [[ "$active_address" == 0x* ]]; then
            curl -v --location --request POST 'https://faucet.testnet.sui.io/gas' \
            --header 'Content-Type: application/json' \
            --data-raw "{
                \"FixedAmountRequest\": {
                    \"recipient\": \"$active_address\"
                }
            }"
        else
            echo "Error: Active address does not start with '0x'."
            exit 1
        fi
    else
        echo "Error: Active environment is not 'testnet'."
        exit 1
    fi
}

if [[ ! -r ~/.sui/sui_config/client.yaml ]]; then
    init_sui
fi

sui client active-address
sui client gas

$*
