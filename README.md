This is component of distributed key/value storage system implemented in Go,
uses AWS SimpleDB for data storage.  Multiple instances of this application
running on independent AWS EC2 instances providing public facing service
through ELB, makes the system resilient to loss of one/more of the EC2
instances, providing high availability and load balancing.  User data
isolation is achieved by storing user data in separate SimpleDB domains. 
SimpleDB backend provides highly available data storage.
