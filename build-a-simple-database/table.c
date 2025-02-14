#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "constants.h"
#include "table.h"
#include "pager.h"

void print_row(Row * row)
{
    printf("(%d, %s, %s)\n", row->id, row->username, row->email);
}

// Write a row to a page.
void serialize_row(Row * source, void *destination)
{
    memcpy(destination + ID_OFFSET, &(source->id), ID_SIZE);
    memcpy(destination + USERNAME_OFFSET, &(source->username), USERNAME_SIZE);
    memcpy(destination + EMAIL_OFFSET, &(source->email), EMAIL_SIZE);
}

// Retrieve a row from a page.
void deserialize_row(void *source, Row * destination)
{
    memcpy(&(destination->id), source + ID_OFFSET, ID_SIZE);
    memcpy(&(destination->username), source + USERNAME_OFFSET, USERNAME_SIZE);
    memcpy(&(destination->email), source + EMAIL_OFFSET, EMAIL_SIZE);
}
