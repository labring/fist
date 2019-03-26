# deploy ldap docker

## start ldap docker

sudo docker run --name my-openldap-container -p 389:389 -p 636:636 -v /root/ldap:/root/ldap   --env LDAP_ORGANISATION="Sealyun Company" --env LDAP_DOMAIN="sealyun.com" --env LDAP_ADMIN_PASSWORD="admin" --detach osixia/openldap:1.2.4



##  into docker 

sudo docker exec -it my-openldap-container /bin/bash

## init  organizational 

cat << EOF > base.ldif
dn: ou=people,dc=sealyun,dc=com
objectClass: organizationalUnit
objectClass: top
ou: people

dn: ou=groups,dc=sealyun,dc=com
objectClass: organizationalUnit
objectClass: top
ou: groups 
EOF

ldapadd -x -w "admin" -D "cn=admin,dc=sealyun,dc=com" -f base.ldif

##  init group 

cat << EOF >group.ldif 
dn: cn=sealyun,ou=groups,dc=sealyun,dc=com
objectClass: top
objectClass: posixGroup
gidNumber: 678
EOF

cat << EOF >group1.ldif 
dn: cn=devops,ou=groups,dc=sealyun,dc=com
objectClass: top
objectClass: posixGroup
gidNumber: 679
EOF

ldapadd -x -w "admin" -D "cn=admin,dc=sealyun,dc=com" -f group.ldif 
ldapadd -x -w "admin" -D "cn=admin,dc=sealyun,dc=com" -f group1.ldif 

## init user  

cat << EOF > adduser.ldif
dn: uid=fanux,ou=people,dc=sealyun,dc=com
objectClass: inetOrgPerson
uid: fanux
sn: 大
cn: 群主
displayName: 大群主
userPassword: fanux
mail: fanux@sealyun.com
EOF

ldapadd -x -w "admin" -D "cn=admin,dc=sealyun,dc=com" -f adduser.ldif 

## add user to group 

cat << EOF >adduser2dev.ldif 
dn: cn=devops,ou=groups,dc=sealyun,dc=com
changetype: modify
add: memberuid
memberuid: fanux
EOF

ldapmodify -x -W -D "cn=admin,dc=sealyun,dc=com" -f adduser2dev.ldif 

cat << EOF >adduser2sealyun.ldif 
dn: cn=sealyun,ou=groups,dc=sealyun,dc=com
changetype: modify
add: memberuid
memberuid: fanux
EOF
ldapmodify -x -W -D "cn=admin,dc=sealyun,dc=com" -f adduser2sealyun.ldif 
