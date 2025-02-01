use wasm_bindgen::prelude::*;

#[derive(Clone, Copy)]
enum Direction {
    Up,
    Right,
    Down,
    Left,
}

#[wasm_bindgen]
pub struct Ant {
    x: usize,
    y: usize,
    direction: Direction,
    grid: Vec<u8>, 
    size: usize,   
}

#[wasm_bindgen]
impl Ant {
    #[wasm_bindgen(constructor)]
    pub fn new(size: usize) -> Ant {
        Ant {
            x: size / 2,
            y: size / 2,
            direction: Direction::Up,
            grid: vec![0; size * size], // 初始化为全白网格
            size,
        }
    }

    pub fn step(&mut self) -> Vec<u8> {
        let index = self.y * self.size + self.x;

        self.grid[index] = 1 - self.grid[index]; 

        self.direction = match (self.grid[index], self.direction) {
            (0, Direction::Up) => Direction::Right,
            (0, Direction::Right) => Direction::Down,
            (0, Direction::Down) => Direction::Left,
            (0, Direction::Left) => Direction::Up,
            (1, Direction::Up) => Direction::Left,
            (1, Direction::Right) => Direction::Up,
            (1, Direction::Down) => Direction::Right,
            (1, Direction::Left) => Direction::Down,
            _ => self.direction,
        };

        match self.direction {
            Direction::Up => {
                if self.y > 0 {
                    self.y -= 1;
                }
            }
            Direction::Right => {
                if self.x < self.size - 1 {
                    self.x += 1;
                }
            }
            Direction::Down => {
                if self.y < self.size - 1 {
                    self.y += 1;
                }
            }
            Direction::Left => {
                if self.x > 0 {
                    self.x -= 1;
                }
            }
        };

        self.grid.clone()
    }

    pub fn x(&self) -> usize {
        self.x
    }

    pub fn y(&self) -> usize {
        self.y
    }
}