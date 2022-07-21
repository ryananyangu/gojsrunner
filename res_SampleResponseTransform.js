const main = (payload) => {

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
