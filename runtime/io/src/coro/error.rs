#[derive(Debug)]
pub struct Error {
    description: String,
}

impl Error {
    pub fn new(description: String) -> Error {
        Error { description }
    }
}

impl std::fmt::Display for Error {
    fn fmt(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(formatter, "{}", self.description)
    }
}

impl std::error::Error for Error {
    fn description(&self) -> &str {
        &self.description
    }
}
