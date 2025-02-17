#include <stdio.h>
#define CHECK(pred) printf("%s...%s\n", #pred, pred? "passed" : "failed")

int main() {
    CHECK(1 == 1);
    CHECK(1 == 2);
    return 0;
}