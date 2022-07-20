const main = (payload) => {
  payload = {
    Envelope: {
      "-xmlns": "http://schemas.xmlsoap.org/soap/envelope/",
      Body: {
        SubmitSMResponse: {
          "-encodingStyle": "http://schemas.xmlsoap.org/soap/encoding/",
          "-soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
          "-xmlns": "http://www.openmindnetworks.com/SoS",
          smResponse: {
            commandStatus: "0",
            messageId: "02940001",
            statusDescription: "success",
            tlvData: "",
            tpTrxID: "3037981",
            transactionID: "HJ87678JHS-01",
          },
        },
      },
      Header: "",
    },
  };
  //TODO: Do your transformation for response here
  return {
    transactionId: Envelope.Body.smResponse.messageId,
    statusId: "",
    date: "",
    statusCode: 200,
    statusDescription: Envelope.Body.smResponse.messageId,
    metadata: {},
  };
};
