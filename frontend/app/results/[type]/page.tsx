"use client";

import { useEffect, useState, use } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { getParticipantScores } from "@/lib/api";
import { getParticipantFromSession } from "@/lib/participant";
import { Score } from "@/types";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";
import { Breadcrumb } from "@/components/ui/Breadcrumb";
import { RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar, ResponsiveContainer, Tooltip } from "recharts";

// ── Scoring helpers ───────────────────────────────────────────────────────────

const SRQ_INDICATORS = [
  { key: "substance_use", label: "Penggunaan Zat (NAPZA)" },
  { key: "psychotic", label: "Gejala Psikotik" },
  { key: "ptsd", label: "Indikasi PTSD" },
] as const;

const IPIP_DIMS = [
  { key: "extraversion", label: "Extraversion", color: "#73fbd3" },
  { key: "agreeableness", label: "Agreeableness", color: "#44e5e7" },
  { key: "conscientiousness", label: "Conscientiousness", color: "#59d2fe" },
  { key: "emotional_stability", label: "Emotional Stability", color: "#4a8fe7" },
  { key: "intellect", label: "Intellect / Openness", color: "#5c7aff" },
] as const;

const API_BASE = "/api";

// ── Component ─────────────────────────────────────────────────────────────────

export default function ResultsPage({ params }: { params: Promise<{ type: string }> }) {
  const router = useRouter();
  const [score, setScore] = useState<Score | null>(null);
  const [participantId, setParticipantId] = useState<string | null>(null);
  const [participantName, setParticipantName] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const resolvedParams = use(params);
  const type = resolvedParams.type; // "srq29" | "ipip"

  useEffect(() => {
    async function loadData() {
      const participant = getParticipantFromSession();
      if (!participant?.participantId) { router.push("/"); return; }

      setParticipantId(participant.participantId);
      setParticipantName(participant.name ?? "");

      try {
        const data = await getParticipantScores(participant.participantId);
        setScore(data);
      } catch {
        setError(true);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, [router]);

  // ─── Loading ─────────────────────────────────────────────────────────────
  if (loading) {
    return (
      <main className="min-h-screen bg-bgLight flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <div className="w-10 h-10 border-4 border-palette4 border-t-transparent rounded-full animate-spin" />
          <p className="text-slate-500 font-semibold text-sm">Memuat hasil...</p>
        </div>
      </main>
    );
  }

  if (error || !score) {
    return (
      <main className="min-h-screen bg-bgLight flex items-center justify-center p-4">
        <BentoCard className="max-w-md w-full text-center">
          <div className="text-4xl mb-4">⚠️</div>
          <h2 className="text-xl font-bold text-textMain mb-2">Hasil tidak ditemukan</h2>
          <p className="text-slate-500 mb-6 font-medium">Mohon selesaikan kuesioner terlebih dahulu.</p>
          <Button variant="primary" onClick={() => router.push("/questionnaire-select")}>Kembali ke Pilihan</Button>
        </BentoCard>
      </main>
    );
  }

  const srq = score.srq_score;
  const ipip = score.ipip_score;

  // IPIP Radar data
  const radarData = IPIP_DIMS.map((d) => ({
    subject: d.label.split(" ")[0],
    value: ipip ? ipip[d.key] : 0,
    fullMark: 50,
  }));

  const exportUrl = `${API_BASE}/export/${participantId}`;

  return (
    <main className="min-h-screen bg-bgLight relative overflow-hidden py-10 px-4">
      <AnimatedBackground />

      <motion.div
        className="max-w-4xl mx-auto relative z-10 space-y-6"
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
      >
        <Breadcrumb items={[
          { label: "Home", href: "/" },
          { label: "Pilih Kuesioner", href: "/questionnaire-select" },
          { label: "Hasil" },
        ]} />

        {/* Hero Card */}
        <BentoCard>
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <span className="text-xs font-black tracking-widest uppercase bg-palette4/10 text-palette4 px-3 py-1.5 rounded-full inline-block mb-3">
                Hasil Interpretasi — {type === "srq29" ? "SRQ-29" : "IPIP-BFM-50"}
              </span>
              <h1 className="text-3xl font-extrabold text-textMain tracking-tight">
                Halo, {participantName} 👋
              </h1>
              <p className="text-slate-500 font-medium mt-1">
                Berikut adalah analisis psikologi Anda berdasarkan jawaban kuesioner.
              </p>
            </div>
            <a
              href={exportUrl}
              download
              className="flex items-center gap-2 px-5 py-2.5 bg-white border-2 border-slate-100 text-slate-600 hover:border-palette4 hover:text-palette4 hover:bg-palette4/5 font-bold rounded-xl transition-all text-sm shadow-sm whitespace-nowrap"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Unduh CSV
            </a>
          </div>
        </BentoCard>

        {/* ── SRQ-29 Results ─────────────────────────────────────────────── */}
        {type === "srq29" && srq && (
          <div className="grid md:grid-cols-2 gap-6">
            {/* Score card */}
            <BentoCard>
              <h2 className="text-sm font-black tracking-wider uppercase text-slate-400 mb-4">Skor Neurotic</h2>
              <div className="flex items-end gap-2 mb-4">
                <span className="text-6xl font-black text-textMain">{srq.neurotic_score}</span>
                <span className="text-slate-300 text-2xl font-bold pb-2">/ 20</span>
              </div>
              <div className={`flex items-center gap-2 px-4 py-3 rounded-xl font-bold text-sm ${srq.neurotic_score >= 6 ? "bg-red-50 text-red-600 border border-red-100" : srq.neurotic_score >= 5 ? "bg-amber-50 text-amber-600 border border-amber-100" : "bg-emerald-50 text-emerald-600 border border-emerald-100"}`}>
                <span className={`w-2.5 h-2.5 rounded-full ${srq.neurotic_score >= 6 ? "bg-red-500" : srq.neurotic_score >= 5 ? "bg-amber-500" : "bg-emerald-500"}`} />
                {srq.neurotic_status}
              </div>

              {/* Mini bar */}
              <div className="mt-5">
                <div className="w-full bg-slate-100 rounded-full h-3 overflow-hidden">
                  <motion.div
                    className={`h-full rounded-full ${srq.neurotic_score >= 6 ? "bg-red-400" : srq.neurotic_score >= 5 ? "bg-amber-400" : "bg-emerald-400"}`}
                    initial={{ width: 0 }}
                    animate={{ width: `${(srq.neurotic_score / 20) * 100}%` }}
                    transition={{ duration: 1, ease: "easeOut", delay: 0.3 }}
                  />
                </div>
                <div className="flex justify-between text-xs text-slate-400 mt-1.5 font-medium">
                  <span>0 — Normal</span>
                  <span>5 — Indikasi</span>
                  <span>20</span>
                </div>
              </div>
            </BentoCard>

            {/* Flags */}
            <BentoCard>
              <h2 className="text-sm font-black tracking-wider uppercase text-slate-400 mb-4">Indikator Tambahan</h2>
              <div className="space-y-3">
                {SRQ_INDICATORS.map(({ key, label }) => {
                  const isPositive = srq[key as keyof typeof srq] === true;
                  return (
                    <div key={key} className={`flex items-center justify-between p-4 rounded-2xl border ${isPositive ? "bg-red-50 border-red-100" : "bg-slate-50 border-slate-100"}`}>
                      <span className={`font-bold text-sm ${isPositive ? "text-red-700" : "text-slate-600"}`}>{label}</span>
                      <span className={`text-xs font-black px-3 py-1 rounded-lg ${isPositive ? "bg-red-100 text-red-600" : "bg-white text-slate-400 border border-slate-200"}`}>
                        {isPositive ? "Ya" : "Tidak"}
                      </span>
                    </div>
                  );
                })}
              </div>
            </BentoCard>
          </div>
        )}

        {/* ── IPIP Results ───────────────────────────────────────────────── */}
        {type === "ipip" && ipip && (
          <>
            {/* Radar Chart */}
            <BentoCard>
              <h2 className="text-sm font-black tracking-wider uppercase text-slate-400 mb-6">Profil Lima Besar Kepribadian</h2>
              <div className="h-64 select-none">
                <ResponsiveContainer width="100%" height="100%">
                  <RadarChart cx="50%" cy="50%" outerRadius="75%" data={radarData}>
                    <PolarGrid stroke="#e2e8f0" />
                    <PolarAngleAxis dataKey="subject" tick={{ fill: "#64748b", fontSize: 12, fontWeight: 700 }} />
                    <PolarRadiusAxis angle={30} domain={[0, 50]} tick={{ fill: "#94a3b8", fontSize: 10 }} />
                    <Radar name="Skor" dataKey="value" stroke="#4a8fe7" fill="#59d2fe" fillOpacity={0.45} strokeWidth={2} />
                    <Tooltip
                      contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", boxShadow: "0 4px 20px rgba(0,0,0,0.08)" }}
                      labelStyle={{ fontWeight: 700, color: "#1e293b" }}
                    />
                  </RadarChart>
                </ResponsiveContainer>
              </div>
            </BentoCard>

            {/* Dimension scores */}
            <div className="grid sm:grid-cols-2 gap-4">
              {IPIP_DIMS.map((dim, i) => {
                const val = ipip[dim.key];
                const pct = (val / 50) * 100;
                return (
                  <motion.div
                    key={dim.key}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: i * 0.07 }}
                    className="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm"
                  >
                    <div className="flex justify-between items-end mb-3">
                      <span className="font-bold text-slate-600 text-sm">{dim.label}</span>
                      <span className="font-black text-2xl text-textMain">{val} <span className="text-xs text-slate-300 font-medium">/ 50</span></span>
                    </div>
                    <div className="w-full bg-slate-100 rounded-full h-2.5 overflow-hidden">
                      <motion.div
                        className="h-full rounded-full"
                        style={{ backgroundColor: dim.color }}
                        initial={{ width: 0 }}
                        animate={{ width: `${pct}%` }}
                        transition={{ duration: 0.8, ease: "easeOut", delay: 0.2 + i * 0.07 }}
                      />
                    </div>
                  </motion.div>
                );
              })}
            </div>
          </>
        )}

        {/* CTA */}
        <BentoCard className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div>
            <p className="font-bold text-textMain">Selesai! Mau lanjut kuesioner lainnya?</p>
            <p className="text-slate-500 text-sm font-medium mt-0.5">Hasil kedua kuesioner dapat diekspor bersama.</p>
          </div>
          <Button variant="primary" onClick={() => router.push("/questionnaire-select")} className="whitespace-nowrap">
            Pilih Kuesioner Lain
          </Button>
        </BentoCard>
      </motion.div>
    </main>
  );
}
