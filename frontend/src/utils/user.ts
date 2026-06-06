export const getUserDisplayName = (userId: number, currentUserId?: number): string => {
  if (currentUserId && userId === currentUserId) {
    return 'You';
  }
  const mockNames: Record<number, string> = {
    1: 'Alice',
    2: 'Bob',
    3: 'Charlie',
    10: 'Alice',
    11: 'Bob',
  };
  return mockNames[userId] || `User ${userId}`;
};
