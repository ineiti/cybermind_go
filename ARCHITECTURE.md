# Architecture

The architecture of cybermind is the following:

```ascii
           User Interface
------------------------------------
 FileSource  |  MailSource  |  Sync
------------------------------------
        Hooks / Authentication
------------------------------------
               CyMiDB
```

## CyberMindDB - CyMiDB

The database of cybermind is set up to support as many links between different nodes as possible. Every node has the 
same structure, but the nodes are separated into different types.

The node structure is as follows:

```go
package main

type CMNode struct {
	ID [32]byte
}
```

Every node is one of the following:
- Device - representing a running instance of cybermind. One cybermind instance can handle more than one user at the 
same time, so that a laptop with multiple user accounts uses only one cybermind instance
- ID - representing a single user. If the user is using more than one device, his ID on all devices will always stay 
the same.
- Hooks - define interactions with external modules. The first modules will be written in go, but other modules might
 be written in other languages and will interact through a REST interface.
- Access Control - TODO - not sure yet how this will look. Useful in two cases: restricting access to user data 
(mostly for syncing), and shared data
- Blob - represents one data in the tree, where the data itself can be stored outside of the database itself 
(filesystem, google drive, ...)
- Links - the core of the cybermind ideas, linking data blobs between each other, adding tags, keywords, search 
terms, ...

### Timeline

One special feature of the CyMiDB is that it has a timeline of all operations, that makes it easy for the syncer 
hooks to check what data did change between two different devices and what needs to be updated. The timeline can be 
queried for changes in the tree.

## Device

The Device node represents one physical device: a computer, a mobile phone, a server.
A database can contain multiple devices, but only one of these can be an 'active Device'.
Devices can also be synchronised across databases.

## Hooks

The first hooks in CyMiDB will be the following:
- File - to add the files of one computer
- Sync - to sync the data to another computer

Hooks are linked to one or more devices where they run.
Only hooks that are linked to the active device will be active.

## UI

The beginning UI will be very simple and only for a local user. Only later versions will have a UI that has 
authentication methods to allow looking at the data when using the internet.
