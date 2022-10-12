# tukxdw
Implementation of  IHE XDW Content Creator, IHE XDW Content Consumer and IHE XDW Content Updator Actors

Requires external Services:-
  IHE DSUB Broker WSE
  IHE PIXm or PIXv3 or PDQv3 server WSE
  MySql DB - connection made using DSN or if url provided by url.

The IHE XDW profile utilises the WS OASIS Human Task standard for the definition of Cross Document Workflow owners, tasks, inputs, outputs, completion behaviours, etc.

The current implementation of tukxdw supports the registering of a XDW definition with the TUK Event Service. The registering process creates DSUB Broker Subscriptions for each XDW input and output task that has a type of '$XDSDocumentEntryTypeCode' in the XDW definition. The resulting broker reference, NHS ID, XDW pathway, topic and expression for each subscription is persisted in the tuk event 'subscriptions' TUK Event service subscriptions DB table. This enables received notifications from a DSUB Broker to be matched to a specific pathway and a specific task in that pathway.

For an example implementation of a DSUB Broker Event Consumer that receives DSUB Broker Notify messages, parses IHE DSUB Notify message and persists the meta data to the TUK Event Service database table 'events', refer to github.com/ipthomas/tukdsub for local deployment and github/ipthomas/tukdsub_lambda for AWS deployment.

To query/update the TUK Event Service Database directly import github.com/ipthomas/tukdbint into your Go project
