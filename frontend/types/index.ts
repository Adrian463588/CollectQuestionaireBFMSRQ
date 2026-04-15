// Types for the questionnaire system

export interface Participant {
  id: string;
  name: string;
  age: number;
  gender: 'male' | 'female';
  created_at: string;
}

export interface ParticipantRequest {
  name: string;
  age: number;
  gender: 'male' | 'female';
}

export interface SRQAnswer {
  [questionNumber: string]: boolean;
}

export interface IPIPAnswer {
  [questionNumber: string]: number;
}

export interface SRQScore {
  neurotic_score: number;
  neurotic_status: string;
  substance_use: boolean;
  psychotic: boolean;
  ptsd: boolean;
}

export interface IPIPScore {
  extraversion: number;
  agreeableness: number;
  conscientiousness: number;
  emotional_stability: number;
  intellect: number;
}

export interface Score {
  id: string;
  participant_id: string;
  srq_score?: SRQScore;
  ipip_score?: IPIPScore;
  created_at: string;
}

export interface SubmissionRequest {
  participant_id: string;
  answers: SRQAnswer | IPIPAnswer;
}

export interface SRQQuestion {
  id: number;
  text: string;
}

export interface IPIPQuestion {
  id: number;
  text: string;
}
