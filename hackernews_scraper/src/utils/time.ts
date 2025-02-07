export const createTimeUtil = () => {
  // Returns a formatted date (YYYY-MM-DD) with the given offset (in days)
  const getFormattedDate = (offset: number = 0): string => {
    const date = new Date();
    date.setDate(date.getDate() + offset);
    return date.toISOString().split('T')[0];
  };

  return {
    today: getFormattedDate(0),
    yesterday: getFormattedDate(-1),
    dayBeforeYesterday: getFormattedDate(-2),
    getDate: getFormattedDate,
  };
};

export type TimeUtil = ReturnType<typeof createTimeUtil>;
