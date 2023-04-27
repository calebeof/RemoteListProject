# Remote List Project
This project aims to develop a Distributed System that consists of a component called RemoteList, which manages a set of lists of integer values, and a set of clients that use the service offered by RemoteList (i.e., insertion, querying, and removal of elements in lists). In this way, RemoteList acts as a server and stores data submitted by clients in lists (more than one list can exist, and each list must have a unique identifier). It also allows remote clients to query data from the list from any position, as well as remove data from the list, but only the value at the end of the list. Multiple clients can use the services offered by RemoteList simultaneously and use common lists.

A synchronous (the client receives confirmation of the operation performed) and persistent (data remains stored even if clients or the server stop executing) communication scheme was implemented using Remote Procedure Call (RPC). In this way, the following operations are available to clients via RPC:

* Append(list_id, v) -> adds the value v at the end of the list with the identifier list_id.
* Get(list_id, i) -> returns the value at position i in the list with the identifier list_id.
* Remove(list_id) -> removes and returns the last element of the list with the identifier list_id.
* Size(list_id) -> obtains the number of elements stored in the list with the identifier list_id.

The system also allows concurrent access to the lists, and ensure data consistency and reliability in the face of failures.

This project is a fork of https://github.com/ruandg/SD_PPGTI and is part of a requirement for the PPGTI Distributed Systems classes.