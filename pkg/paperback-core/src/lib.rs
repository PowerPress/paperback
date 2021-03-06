/*
 * paperback: paper backup generator suitable for long-term storage
 * Copyright (C) 2018-2020 Aleksa Sarai <cyphar@cyphar.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

#![forbid(unsafe_code)]

extern crate aead;
extern crate bip39;
extern crate chacha20poly1305;
extern crate ed25519_dalek;
extern crate itertools;
extern crate nom;
extern crate rand;
extern crate serde;
extern crate unsigned_varint;
extern crate zbase32;

#[cfg(test)]
extern crate quickcheck;
#[cfg(test)]
#[macro_use]
extern crate quickcheck_macros;

/// Implementation of Shamir Secret Sharing.
mod shamir;

/// `nom` helpers which haven't been upstreamed to the relevant projects.
mod nom_helpers;

/// Initial version of paperback wire format types.
///
/// This module also includes all of the necessary code to serialise and
/// interact with the relevant structures.
pub mod v0;

/// Re-export of the newest paperback wire format types.
pub use v0 as latest;
