mod utils;

use wasm_bindgen::prelude::*;

#[wasm_bindgen]
#[derive(Clone, Copy, Debug, PartialEq)]
#[repr(u8)]
pub enum Cell {
    Dead = 0,
    Alive = 1,
}

#[wasm_bindgen]
impl Cell {
    pub fn toggle(&mut self) {
        *self = match *self {
            Cell::Dead => Cell::Alive,
            Cell::Alive => Cell::Dead,
        }
    }
}

#[wasm_bindgen]
pub struct Universe {
    width: u32,
    height: u32,
    cells: Vec<Cell>,
}

#[wasm_bindgen]
impl Universe {
    fn get_index (&self, row: u32, column: u32) -> usize {
        (row * self.width + column) as usize
    }

    fn count_live_neighbors(&self, row: u32, column: u32) -> u8 {
        let mut count = 0;

        for i in [self.height - 1, 0, 1] {
            for j in [self.width - 1, 0, 1] {
                if i == 0 && j == 0 {
                    continue;
                }

                let r = (row + i) % self.height;
                let c = (column + j) % self.width;
                let idx = self.get_index(r, c);
                count += self.cells[idx] as u8;
            }
        }
        count
    }
    pub fn tick(&mut self) {
        let mut cells = self.cells.clone();
        for r in 0..self.height {
            for c in 0..self.width {
                let idx = self.get_index(r, c);
                cells[idx] = match (self.cells[idx], self.count_live_neighbors(r,c)) {
                    (_, 3) => Cell::Alive,
                    (Cell::Alive, 2) => Cell::Alive,
                    _ => Cell::Dead,
                }
            }
        }
        self.cells = cells;
    }
    pub fn render(&self) -> String {
        self.to_string()
    }
    pub fn new() -> Self {
        let width = 64;
        let height = 64;
        let cells = (0.. width *height ).map(|x|{
            if x % 2 == 0 || x % 7 == 0 {
                Cell::Alive
            } else {
                Cell::Dead
            }
        }).collect();
        Self {
            width,
            height,
            cells,
        }
    }
    pub fn width(&self) -> u32 {
        self.width
    }
    pub fn height(&self) -> u32 {
        self.height
    }
    pub fn cells(&self) -> *const Cell {
        self.cells.as_ptr()
    }
    pub fn toggle_cell(&mut self, row: u32, column: u32) {
        let idx = self.get_index(row, column);
        self.cells[idx].toggle();
    }
}

use std::fmt;

impl fmt::Display for Universe {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        for r in 0..self.height {
            for c in 0..self.width {
                let idx = self.get_index(r, c);
                let x = if self.cells[idx] == Cell::Dead{'□'} else {'■'};
                write!(f, "{}", x)?;
            }
            write!(f, "\n")?;
        }
        Ok(())
    }
}

