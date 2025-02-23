use crate::config::config::AppConfig;
use sqlx::MySqlPool;

use std::collections::HashSet;
use std::fs;
use std::io::{BufRead, BufReader, Write};
use std::path::Path;

const PUBLISHED_FILE: &str = "published_articles.txt";
const MAX_PUBLISHED_LINES: usize = 1000;

/// Loads the set of published article links from a file.
/// If the file does not exist, returns an empty set.
pub fn load_published_articles() -> HashSet<String> {
    let mut published = HashSet::new();
    if Path::new(PUBLISHED_FILE).exists() {
        if let Ok(file) = fs::File::open(PUBLISHED_FILE) {
            let reader = BufReader::new(file);
            for line in reader.lines().flatten() {
                published.insert(line);
            }
        }
    }
    published
}

/// Prunes the published articles set if it reaches MAX_PUBLISHED_LINES.
pub fn prune_published_articles(published: &mut HashSet<String>) {
    if published.len() >= MAX_PUBLISHED_LINES {
        published.clear();
    }
}

/// Updates the published articles file by writing all the links from the set.
/// This helper prunes the set if needed, then recreates the file.
pub fn update_published_articles(published: &HashSet<String>) {
    let mut pruned = published.clone();
    prune_published_articles(&mut pruned);

    if let Ok(mut file) = fs::File::create(PUBLISHED_FILE) {
        for link in pruned {
            let _ = writeln!(file, "{}", link);
        }
    }
}

/// Creates and returns a MySQL connection pool.
pub async fn create_pool(config: &AppConfig) -> Result<MySqlPool, sqlx::Error> {
    MySqlPool::connect(&config.database_url).await
}
