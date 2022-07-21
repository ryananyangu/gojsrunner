const main = (payload, headers, constants) => {
  const {
    MSISDN,
    accountNumber,
    transactionId,
    amount,
    currentDate,
    narration,
    ISOCurrencyCode,
    customerName,
    paymentMode,
    callback,
    metadata,
  } = payload;

  // TODO: //Put your transformation code here
  return {
    requeststring: "",
    headers: headers,
  };
};

// You can add extra functions here so long as they are invoked in the main function
