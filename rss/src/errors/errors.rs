use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use std::net::AddrParseError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),

    #[error("Article fetch error: {0}")]
    ArticleFetch(String),

    #[error("RSS build error: {0}")]
    RssBuild(String),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Address parse error: {0}")]
    AddrParse(#[from] AddrParseError),

    #[error("Service error: {0}")]
    BoxedError(#[from] Box<dyn std::error::Error + Send + Sync>),
}

impl From<axum::http::Error> for AppError {
    fn from(err: axum::http::Error) -> Self {
        AppError::BoxedError(Box::new(err))
    }
}

// Update IntoResponse implementation
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let status = match &self {
            AppError::Database(_) => StatusCode::INTERNAL_SERVER_ERROR,
            AppError::ArticleFetch(_) => StatusCode::BAD_GATEWAY,
            AppError::RssBuild(_) => StatusCode::INTERNAL_SERVER_ERROR,
            AppError::Io(_) => StatusCode::INTERNAL_SERVER_ERROR,
            AppError::AddrParse(_) => StatusCode::BAD_REQUEST,
            AppError::BoxedError(_) => StatusCode::INTERNAL_SERVER_ERROR,
        };

        (status, self.to_string()).into_response()
    }
}
