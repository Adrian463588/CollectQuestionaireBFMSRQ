import { ParticipantRequest } from "@/types";

type StoredParticipant = ParticipantRequest & {
  participantId?: string;
};

const SESSION_KEY = "participant";

export function getParticipantFromSession(): StoredParticipant | null {
  const raw = sessionStorage.getItem(SESSION_KEY);
  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw) as StoredParticipant;
  } catch {
    return null;
  }
}

export function saveParticipantToSession(participant: StoredParticipant): void {
  sessionStorage.setItem(SESSION_KEY, JSON.stringify(participant));
}

export function saveParticipantIdToSession(participantId: string): void {
  const participant = getParticipantFromSession();
  if (!participant) {
    return;
  }

  saveParticipantToSession({
    ...participant,
    participantId,
  });
}
