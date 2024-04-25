// processSummary function processes a summary object and returns a string.
// If the summary is valid and not empty, it returns the summary; otherwise, it returns a default message.
export const processSummary = (
  summary: { String: string; Valid: boolean } | null,
): string => {
  return summary?.Valid && summary.String.trim() !== ''
    ? summary.String
    : 'No summary available';
};

// formatDate function takes a date string and formats it to a human-readable date.
// If the input date string is invalid or empty, it returns appropriate messages.
export const formatDate = (dateStr: string): string => {
  if (!dateStr) {
    return 'Date not available';
  }
  const date = new Date(dateStr);
  return isNaN(date.getTime())
    ? 'Invalid Date'
    : date.toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    });
};
