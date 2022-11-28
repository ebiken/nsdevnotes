#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <sys/types.h>
#include <sys/socket.h>

#include <unistd.h>

#include <linux/netlink.h>
#include <linux/rtnetlink.h>

int main(int argc, char **argv) {
    
    int fd; // file descripter for netlink socket
    struct sockaddr_nl local;

    pid_t pid = getpid();

    fd = socket(AF_NETLINK, SOCK_RAW, NETLINK_ROUTE);

    memset(&local, 0, sizeof(local)); /* fill-in local address information */
    local.nl_family = AF_NETLINK;
    local.nl_pid = pid;
    local.nl_groups = RTMGRP_IPV4_ROUTE;

    if (bind(fd, (struct sockaddr *) &local, sizeof(local)) < 0) {
        // cannot bind socket
        return -1;
    }

    return 0;
}
