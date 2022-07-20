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

  let requeststring = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://topupretail.com/">
<SOAP-ENV:Body>
<ns1:AquaPayment>
<ns1:req>
<ns1:terminalMsgID>${transactionId}</ns1:terminalMsgID>
<ns1:terminalID>${constants.terminalID}</ns1:terminalID>
<ns1:msgID>${constants.terminalID}</ns1:msgID>
<ns1:authCred>
<ns1:opName>${constants.opName}</ns1:opName>
<ns1:password>${constants.password}</ns1:password>
</ns1:authCred>
<ns1:clientNumber>${accountNumber}</ns1:clientNumber>
<ns1:entityNumber>${metadata.entityNumber}</ns1:entityNumber>
<ns1:documentTypeID>${metadata.documentTypeID}</ns1:documentTypeID>
<ns1:documentYear>${metadata.documentYear}</ns1:documentYear>
<ns1:documentNumber>${metadata.documentNumber}</ns1:documentNumber>
<ns1:clientName>${customerName}</ns1:clientName>
<ns1:purchaseValue>${amount}</ns1:purchaseValue>
<ns1:receiptFormat>${constants.receiptFormat}</ns1:receiptFormat>
<ns1:terminalLocation />
<ns1:terminalChannel>${paymentMode}</ns1:terminalChannel>
<ns1:terminalCompanyName>${constants.terminalCompanyName}</ns1:terminalCompanyName>
<ns1:terminalOperator>${constants.terminalOperator}</ns1:terminalOperator>
</ns1:req>
</ns1:AquaPayment>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`;

  console.log(payload, headers, constants);

  // TODO: //Put your transformation code here
  return {
    requeststring: requeststring.replace(/(\r\t\n|\n|\r)/gm, ""),
    headers: headers,
  };
};

// You can add extra functions here so long as they are invoked in the main function
