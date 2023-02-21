<!-- Notes from class -->

# Final project requirements

* Groups of up to 3
* Private git repo, add @kbsec as a contributor
* There will be milestones we must hit to get full credit

Rough milestones:
* 2/17 - pick teammates
* 2/20 - Schedule roadmap meeting with kai
* 2/20 - execution
* 2/27 - http client, basic team server outline
* 3/6 - basic rpc, pick your special feature
* 3/20 - file io
* 3/27 - Process Injection (shellcode)
* 4/6 - Reflective DLL payload
* more on slides

## Implant
Targets the windows operating system

Requirements:
* RPC and C2 channel
  * Implant must communicate to the C2 server using an RPC framework built on top of a C2 channel and serialization format
  * Examples RPCS: JSON, TLV, SOAP, Protobuf
  * RPC has to be async (need to have long-polling set up, with jitter)
  * C2 channel has to be HTTPS to start, can add another
  * Http comments:
    * Need to use TLS (https)
    * Implant must be proxy aware
    * Implant shouldnt crash if internet cuts out
    * Each implant will be compiled with configuration data including
    	* The C2(s) to communicate back to
	* The amount of time to go dormant before sleeping again
	* Add randomness to the sleep time
	* Kill date of the implant
	* The server's public key
	* Modifying the above at runtime
* Cryptography
  * Must use *secure* cryptography
  * Building blocks: Asymmetric crypto (RSA), Symmetric crypto (AES), Hashing (SHA256)
  * Session establishment: Implant should establish a secure connection to a C2 server, and identify itself
  * Session update: re-establish a secure connection to a C2
  * Recommended authenticated cryptography: AES-GCM
  * Hash functions:
    * You must payload to a device that has the C:\Malware\ch0nky.txt file path using a secure hash function. Can't just check for the file, must check fo the hash of the file name
    * Bonus: Have the stager report back the MAC address or CPU ID or other identifying information and tailor the payload to only run on a machine that it is supposed to
  * Don't roll own crypto, use Windows CNG (or something else like libsodium, libressl, mbedtls. NOT openssl)
* Situation awareness
  * Whereami, whoami, whatami
  * Must be able to do the following:
    * List computers network interfaces (MAC, IPs, interface names... etc)
    * Read environment variables
    * Get the windows version
    * Get the current username and token
    * Get the computers name
    * Get the machine GUID
    * list the files in a directory
    * Change directory
    * List all running processes
    * Bonus: Anti-sandbox detection
    * Devise a way to limit the number of implants running on the same host (mutex, named pipe)
    * Could do tiered execition (eg, 3 max executing binaries, 1 for interactive, 1 for recovery, one for long haul)

* Execution
  * CreateProcess and redirect I/O
  * Chellcode execution: support shellcode execution  in a local process
  * **Seatbelt**
  * Payload must be available in multiple formats
    * EXE
    * DLL
    * Reflective DLL
    * Bonus: Position independent shellcode
* Injection
  * Process injection: support execution of shelcode in either a remote or local process
  * Look up Donut
* File I/O
  * Should be able to:
    * Read files
    * Download files from C2 and write them either encrpyted or unencrypted
    * Interact with the disk (safely)
* Persistence
  * Need to be able to persists on the target past reboots
  * Strive to make all payloads file-less, need to implement a stager
  * Stager: wakes up, beacons to C2, grabs full payload, and executes it in memory
  * Bonus points: tailor the payload to the machine (ie the exe will only run if the machine has the same MAC address, CPU ID, etc)
* Loot
  * Programatically loot information from victim machine
  * Full credit: Chrome passwords and cookies from default user profile
  * Extra credit: All chromium based passwords, cookies, autofill, etc from all user profiles
  * NOT ALLOWED: keylogging and bitcoin wallet stealing
* Defense Evasion
  * Encrypt/obfuscate configuration strings
  * No using powershell for situational awareness tasks
  * Select one of the folowing
    * A crypter that packs and obfuscates your payload, no upx, but can use a modified version of upx ;)
    * A method for defeating AMSI
    * API hashing / dynamic IAT
    * An RPC to mimic a legitimate service
    * Payload tailoring
    * Something else of your choosing
  
## Opsec
Don't get caught
* Minimal writing of files to disk
* Minimal use of cmd/powershell
* **Minimal allocation of r/w/x memory**, can change memory permissions
* Minimal creation of new processes/threads
* Separate low-latency comunications fro mhigh-latency actions
* Keep payloads modular and small. Load functionality you need when you need it
* Be able to explain *why* your implant is doing someting a specific way

## General comments
Staff will support people who:
* Write their C2 Server and client in python
* Write the backent message broker with RabbidMQ, Redis, or ZeroMQ
* The client UI with prompt toolkit
* Implant in C/C++
* For implant:
  * Use C, C++, Zig, Rust

## C2 requirements
* handle connections from multiple operators, and multiple implants
* Use flask with Gunicorn as WSGI
* Use postgres or MySql as the backend database
* Use flask-SQLAlchemy to facilitate CRUD operations on agents and operators
* Using RabbitMQ/ZeroMQ/Redis to broker messages between implant and C2, and client and C2

## C2 -> DB
* Implants
* Commands
* Jobs
* Clients (operators)

## C2 -> Db -> Implants
* Implant ID
* Computer name
* Username
* Guid
* Integrity
* connecting IP addr
* Session key (asym key)
* Sleep
* Jitter
* First seen
* Last seen
* Expected check in 

## Messaging
* Operators should be notified whenever one of the following occurs
* A new implant connects to the C2 for the first time
* A client connects to the C2
* An operator issues a command to an agent
* An agent responds to a job that was issued by an operator
* RabbitMQ Pub sub is probably the easiest way to handle this
* Your client code cna be run on the same box as the C2. 

## Client
* Client needs to be able to securely connect to the C2, send and receive data, and get updates fom the agent
* This can be a terminal interface or web interface

## Special feature
* You must implement an advanced C2 feature
* Examples on slides

## Rules of thumb
* No functions with more than 50 lines of code
* No files with more than 500 lines of code
* Needs to run on the windows VM setup for this class
* All commands need to be documented

## Build recommendations
* Swith with ming2
* Use resource files
