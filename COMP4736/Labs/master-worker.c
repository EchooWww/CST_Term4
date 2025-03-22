#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include <sys/wait.h>
#include <pthread.h>

int item_to_produce, curr_buf_size;
int total_items, max_buf_size, num_workers, num_masters;

int *buffer;

// Synchronization variables
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t full = PTHREAD_COND_INITIALIZER;
pthread_cond_t empty = PTHREAD_COND_INITIALIZER;
int items_consumed = 0;

void print_produced(int num, int master) {
    printf("Produced %d by master %d\n", num, master);
}

void print_consumed(int num, int worker) {
    printf("Consumed %d by worker %d\n", num, worker);
}

// produce items and place in buffer
// modified to synchronize correctly
void *generate_requests_loop(void *data)
{
    int thread_id = *((int *)data);

    while (1) {
        // Lock the mutex before checking if we can produce
        pthread_mutex_lock(&mutex);
        
        // Check if we've reached the total number of items to produce
        if (item_to_produce >= total_items) {
            pthread_mutex_unlock(&mutex);
            break;
        }
        
        // Wait if the buffer is full
        while (curr_buf_size >= max_buf_size) {
            pthread_cond_wait(&empty, &mutex);
        }
        
        // Add item to buffer
        buffer[curr_buf_size++] = item_to_produce;
        print_produced(item_to_produce, thread_id);
        item_to_produce++;
        
        // Signal that the buffer is not empty anymore
        pthread_cond_signal(&full);
        pthread_mutex_unlock(&mutex);
    }
    
    // Once all items are produced, signal any waiting consumers
    pthread_mutex_lock(&mutex);
    pthread_cond_broadcast(&full);
    pthread_mutex_unlock(&mutex);
    
    return NULL;
}

// function to be run by worker threads
void *consume_requests_loop(void *data)
{
    int thread_id = *((int *)data);
    int consumed_item;
    
    while (1) {
        // Lock the mutex before checking if we can consume
        pthread_mutex_lock(&mutex);
        
        // Check if we've consumed all items
        if (items_consumed >= total_items) {
            pthread_mutex_unlock(&mutex);
            break;
        }
        
        // Wait if the buffer is empty
        while (curr_buf_size <= 0 && items_consumed < total_items) {
            pthread_cond_wait(&full, &mutex);
            
            // After waking up, check again if all items are consumed
            if (items_consumed >= total_items) {
                pthread_mutex_unlock(&mutex);
                return NULL;
            }
        }
        
        // All items produced and consumed
        if (curr_buf_size <= 0) {
            pthread_mutex_unlock(&mutex);
            break;
        }
        
        // Consume item from buffer (take from the start of buffer and shift everything left)
        consumed_item = buffer[0];
        for (int i = 0; i < curr_buf_size - 1; i++) {
            buffer[i] = buffer[i + 1];
        }
        curr_buf_size--;
        items_consumed++;
        
        // Print consumption
        print_consumed(consumed_item, thread_id);
        
        // Signal that the buffer is not full anymore
        pthread_cond_signal(&empty);
        pthread_mutex_unlock(&mutex);
    }
    
    return NULL;
}

// main program
int main(int argc, char *argv[])
{
    int *master_thread_id;
    pthread_t *master_thread;
    int *worker_thread_id;
    pthread_t *worker_thread;
    
    item_to_produce = 0;
    curr_buf_size = 0;
    items_consumed = 0;
  
    int i;
  
    if (argc < 5) {
        printf("./master-worker #total_items #max_buf_size #num_workers #num_masters, e.g.\n ./master-worker 10000 1000 4 3\n");
        exit(1);
    }
    else {
        num_masters = atoi(argv[4]);
        num_workers = atoi(argv[3]);
        total_items = atoi(argv[1]);
        max_buf_size = atoi(argv[2]);
    }

    buffer = (int *)malloc(sizeof(int) * max_buf_size);

    // Create master producer threads
    master_thread_id = (int *)malloc(sizeof(int) * num_masters);
    master_thread = (pthread_t *)malloc(sizeof(pthread_t) * num_masters);
    for (i = 0; i < num_masters; i++) {
        master_thread_id[i] = i;
    }

    // Create worker consumer threads
    worker_thread_id = (int *)malloc(sizeof(int) * num_workers);
    worker_thread = (pthread_t *)malloc(sizeof(pthread_t) * num_workers);
    for (i = 0; i < num_workers; i++) {
        worker_thread_id[i] = i;
    }
    
    // Initialize mutex and condition variables
    pthread_mutex_init(&mutex, NULL);
    pthread_cond_init(&full, NULL);
    pthread_cond_init(&empty, NULL);

    // Create the worker threads
    for (i = 0; i < num_workers; i++) {
        pthread_create(&worker_thread[i], NULL, consume_requests_loop, (void *)&worker_thread_id[i]);
    }

    // Create the master threads
    for (i = 0; i < num_masters; i++) {
        pthread_create(&master_thread[i], NULL, generate_requests_loop, (void *)&master_thread_id[i]);
    }

    // Wait for all master threads to complete
    for (i = 0; i < num_masters; i++) {
        pthread_join(master_thread[i], NULL);
        printf("master %d joined\n", i);
    }

    // Wait for all worker threads to complete
    for (i = 0; i < num_workers; i++) {
        pthread_join(worker_thread[i], NULL);
        printf("worker %d joined\n", i);
    }

    // Destroy mutex and condition variables
    pthread_mutex_destroy(&mutex);
    pthread_cond_destroy(&full);
    pthread_cond_destroy(&empty);

    // Deallocating Buffers
    free(buffer);
    free(master_thread_id);
    free(master_thread);
    free(worker_thread_id);
    free(worker_thread);

    return 0;
}