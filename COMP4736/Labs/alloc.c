#include "alloc.h"

/* Memory block metadata structure - kept separate from the memory page */
typedef struct block_meta {
    size_t size;            // Size of the block
    int free;               // 1 if free, 0 if allocated
    size_t offset;          // Offset from the start of the memory page
    struct block_meta *next; // Next block in the list
} block_meta;

/* Global variables */
static void *memory_page = NULL;    // Pointer to the 4KB memory page
static block_meta *meta_list = NULL; // List of memory block metadata

/* Initialize memory manager */
int init_alloc() {
    // Allocate a 4KB page using mmap
    memory_page = mmap(NULL, PAGESIZE, PROT_READ | PROT_WRITE, 
                      MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    
    if (memory_page == MAP_FAILED) {
        return 1; // mmap failed
    }
    
    // Initialize the metadata for the single free block covering the whole page
    meta_list = malloc(sizeof(block_meta));
    if (meta_list == NULL) {
        munmap(memory_page, PAGESIZE);
        return 1; // malloc failed
    }
    
    meta_list->size = PAGESIZE;
    meta_list->free = 1;
    meta_list->offset = 0;
    meta_list->next = NULL;
    
    return 0; // Success
}

/* Clean up memory manager */
int cleanup() {
    // Free all metadata blocks
    block_meta *current = meta_list;
    block_meta *next;
    
    while (current != NULL) {
        next = current->next;
        free(current);
        current = next;
    }
    
    meta_list = NULL;
    
    // Unmap the memory page
    if (munmap(memory_page, PAGESIZE) == -1) {
        return 1; // munmap failed
    }
    
    memory_page = NULL;
    
    return 0; // Success
}

/* Find a suitable free block for allocation (first-fit strategy) */
static block_meta* find_free_block(size_t size) {
    block_meta *current = meta_list;
    
    while (current != NULL) {
        if (current->free && current->size >= size) {
            return current; // Found a suitable block
        }
        current = current->next;
    }
    
    return NULL; // No suitable block found
}

/* Split a block if necessary */
static void split_block(block_meta *block, size_t size) {
    // Only split if the remaining space is at least MINALLOC bytes
    if (block->size - size >= MINALLOC) {
        block_meta *new_block = malloc(sizeof(block_meta));
        if (new_block == NULL) {
            return; // Cannot split, just use the entire block
        }
        
        // Setup the new block
        new_block->size = block->size - size;
        new_block->free = 1;
        new_block->offset = block->offset + size;
        new_block->next = block->next;
        
        // Update the current block
        block->size = size;
        block->next = new_block;
    }
}

/* Merge adjacent free blocks */
static void merge_blocks() {
    block_meta *current = meta_list;
    
    while (current != NULL && current->next != NULL) {
        if (current->free && current->next->free && 
            current->offset + current->size == current->next->offset) {
            // The blocks are adjacent and both free - merge them
            current->size += current->next->size;
            
            // Remove the next block from the list
            block_meta *to_free = current->next;
            current->next = current->next->next;
            free(to_free);
            
            // Don't advance current as we might be able to merge with the next block as well
        } else {
            current = current->next;
        }
    }
}

/* Allocate memory */
char *alloc(int size) {
    // Size must be positive and a multiple of MINALLOC
    if (size <= 0 || size % MINALLOC != 0) {
        return NULL;
    }
    
    // Find a suitable free block
    block_meta *block = find_free_block(size);
    if (block == NULL) {
        return NULL; // No suitable block found
    }
    
    // Split the block if it's larger than needed
    split_block(block, size);
    
    // Mark the block as allocated
    block->free = 0;
    
    // Return a pointer to the memory in the page
    return (char *)memory_page + block->offset;
}

/* Free allocated memory */
void dealloc(char *ptr) {
    if (ptr == NULL || memory_page == NULL) {
        return;
    }
    
    // Calculate the offset from the start of the memory page
    size_t offset = ptr - (char *)memory_page;
    
    // Find the block with this offset
    block_meta *current = meta_list;
    while (current != NULL) {
        if (current->offset == offset) {
            // Mark the block as free
            current->free = 1;
            
            // Try to merge with adjacent blocks
            merge_blocks();
            return;
        }
        current = current->next;
    }
}