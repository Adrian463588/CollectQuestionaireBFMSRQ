import {
  IPIPAnswer,
  Participant,
  ParticipantRequest,
  Score,
  SRQAnswer,
} from "@/types";

export const API_BASE = process.env.NEXT_PUBLIC_API_URL || "/api";

type QuestionnaireType = "srq29" | "ipip-bfm-50";

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });

  if (!response.ok) {
    let errorMessage = `Request failed with status ${response.status}`;
    try {
      const payload = (await response.json()) as { error?: string };
      if (payload.error) errorMessage = payload.error;
    } catch {
      // Keep the default message when response body is not JSON.
    }
    throw new Error(errorMessage);
  }

  return (await response.json()) as T;
}

export async function createParticipant(
  payload: ParticipantRequest,
): Promise<Participant> {
  return apiFetch<Participant>("/participants", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function submitQuestionnaire(
  participantId: string,
  questionnaireType: QuestionnaireType,
  answers: SRQAnswer | IPIPAnswer,
): Promise<void> {
  await apiFetch("/responses", {
    method: "POST",
    body: JSON.stringify({
      participant_id: participantId,
      questionnaire_type: questionnaireType,
      answers,
    }),
  });
}

export async function getParticipantScores(
  participantId: string,
): Promise<Score> {
  return apiFetch<Score>(`/scores/${participantId}`);
}

export async function exportParticipantCSV(
  participantId: string,
): Promise<Blob> {
  const response = await fetch(`${API_BASE}/export/${participantId}`);
  if (!response.ok) throw new Error("Failed to export participant data");
  return response.blob();
}

export type DashboardItem = {
  participant: Participant;
  score: Score | null;
};

export async function fetchAllDashboardData(): Promise<DashboardItem[]> {
  return apiFetch<DashboardItem[]>("/dashboard");
}
