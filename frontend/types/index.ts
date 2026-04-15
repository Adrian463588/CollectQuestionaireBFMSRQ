// Types for the questionnaire system

export interface Participant {
  id: string;
  name: string;
  age: number;
  gender: 'male' | 'female';
  created_at: string;
  deleted_at?: string | null;
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

/**
 * SRQ-29 scoring result.
 * Domain-based per WHO guidelines:
 *   Q1–20  → Gangguan Mental Emosional (GME)
 *   Q21    → Penggunaan Zat Psikoaktif
 *   Q22–24 → Gejala Psikotik
 *   Q25–29 → Gejala PTSD
 */
export interface SRQScore {
  // ── Core domain fields (backward compatible) ────────────────────────────
  neurotic_score: number;    // sum of YA on Q1–20
  neurotic_status: string;   // "normal" | "indikasi_gme" | "rekomendasi_rujukan"
  substance_use: boolean;    // Q21 = YA
  psychotic: boolean;        // any of Q22–24 = YA
  ptsd: boolean;             // any of Q25–29 = YA
  // ── Extended domain detail ───────────────────────────────────────────────
  psychotic_count: number;   // count of YA in Q22–24 (0–3)
  ptsd_count: number;        // count of YA in Q25–29 (0–5)
  // ── Aggregate scores ─────────────────────────────────────────────────────
  total_score: number;       // sum of YA Q1–29
  overall_risk: string;      // "rendah" | "sedang" | "tinggi" | "kritis"
  // ── Research dummy variables (0/1) ────────────────────────────────────────
  emotional_disorder: number;
  substance_dummy: number;
  psychotic_dummy: number;
  ptsd_dummy: number;
}

/**
 * IPIP-BFM-50 scoring result.
 * Scores are MEAN of 10 items per dimension (range 1.0–5.0)
 * after reverse-scoring negatively-keyed items.
 * Reference: Akhtar & Azwar (2019), ipip.ori.org
 */
export interface IPIPScore {
  // ── Mean scores (1.0–5.0) ────────────────────────────────────────────────
  extraversion: number;
  agreeableness: number;
  conscientiousness: number;
  emotional_stability: number;
  intellect: number;
  // ── Interpretation labels ─────────────────────────────────────────────────
  // "sangat_tinggi" | "tinggi" | "rata_rata" | "rendah" | "sangat_rendah"
  extra_label: string;
  agre_label: string;
  cons_label: string;
  stab_label: string;
  intell_label: string;
  // ── Raw sum scores (10–50) ───────────────────────────────────────────────
  extraversion_sum: number;
  agreeableness_sum: number;
  conscientiousness_sum: number;
  emotional_stability_sum: number;
  intellect_sum: number;
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
