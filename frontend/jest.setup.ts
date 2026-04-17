import '@testing-library/jest-dom';

// Mock Next.js router
jest.mock('next/navigation', () => ({
  useRouter() {
    return {
      push: jest.fn(),
      replace: jest.fn(),
      prefetch: jest.fn(),
      back: jest.fn(),
    };
  },
}));

// Mock simple API client logic
jest.mock('@/lib/api', () => ({
  createParticipant: jest.fn().mockResolvedValue({ id: 'mock-uuid', name: 'Test' }),
  submitQuestionnaire: jest.fn().mockResolvedValue(true),
}));

jest.mock('@/lib/participant', () => ({
  getParticipantFromSession: jest.fn().mockReturnValue({ participantId: 'mock-uuid', name: 'Test User' }),
  saveParticipantIdToSession: jest.fn(),
}));
