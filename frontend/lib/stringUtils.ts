export const processSummary = (
    summary: { String: string; Valid: boolean } | null,
  ): string => {
    return summary?.Valid && summary.String.trim() !== ''
      ? summary.String
      : 'No summary available';
};
  
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
