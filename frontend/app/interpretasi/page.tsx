"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import {
  RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar,
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, Cell,
} from "recharts";
import { fetchAllDashboardData, DashboardItem, API_BASE } from "@/lib/api";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";
import { AdminGuard } from "@/components/AdminGuard";
import { Breadcrumb } from "@/components/ui/Breadcrumb";

// ── Constants ─────────────────────────────────────────────────────────────────

const IPIP_DIMS = [
  { key: "extraversion", label: "Extraversion", color: "#73fbd3" },
  { key: "agreeableness", label: "Agreeableness", color: "#44e5e7" },
  { key: "conscientiousness", label: "Conscientiousness", color: "#59d2fe" },
  { key: "emotional_stability", label: "Emotional Stability", color: "#4a8fe7" },
  { key: "intellect", label: "Intellect", color: "#5c7aff" },
] as const;

// ── Helpers ───────────────────────────────────────────────────────────────────

function getSRQClass(s: DashboardItem["score"]): string {
  const srq = s?.srq_score;
  if (!srq) return "-";
  const flags: string[] = [];
  if (srq.neurotic_score >= 5) flags.push("Indikasi GME");
  if (srq.substance_use) flags.push("Penggunaan Zat");
  if (srq.psychotic) flags.push("Gejala Psikotik");
  if (srq.ptsd) flags.push("Gejala PTSD");
  return flags.length > 0 ? flags.join(" · ") : "Normal";
}

function getDominantTrait(s: DashboardItem["score"]): string {
  const ipip = s?.ipip_score;
  if (!ipip) return "-";
  return IPIP_DIMS.reduce((best, dim) => {
    const a = Number(ipip[dim.key] ?? 0);
    const b = Number(ipip[best.key as keyof typeof ipip] ?? 0);
    return a > b ? dim : best;
  }).label;
}

function hasSevereFlag(s: DashboardItem["score"]): boolean {
  return !!s?.srq_score && (s.srq_score.psychotic || s.srq_score.ptsd || s.srq_score.substance_use);
}

// ── Component ─────────────────────────────────────────────────────────────────

export default function InterpretasiPage() {
  const [data, setData] = useState<DashboardItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAllDashboardData()
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-bgLight flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="w-10 h-10 border-4 border-palette4 border-t-transparent rounded-full animate-spin" />
          <p className="text-slate-500 font-semibold text-sm">Memuat data...</p>
        </div>
      </div>
    );
  }

  // ── Aggregated stats ────────────────────────────────────────────────────────
  const withSRQ = data.filter((d) => d.score?.srq_score);
  const withIPIP = data.filter((d) => d.score?.ipip_score);
  const avgNeurotic = withSRQ.length
    ? withSRQ.reduce((acc, d) => acc + (d.score!.srq_score!.neurotic_score ?? 0), 0) / withSRQ.length
    : 0;

  const avgIPIP = IPIP_DIMS.map((dim) => ({
    subject: dim.label.split(" ")[0],
    value: withIPIP.length
      ? withIPIP.reduce((acc, d) => acc + Number(d.score!.ipip_score![dim.key] ?? 0), 0) / withIPIP.length
      : 0,
    fullMark: 5,
    color: dim.color,
  }));

  const srqBarData = withSRQ.map((d) => ({
    name: d.participant.name.split(" ")[0].slice(0, 8),
    skor: d.score!.srq_score!.neurotic_score,
  }));

  const urgentCount = data.filter((d) => hasSevereFlag(d.score)).length;

  // ── Export URLs ─────────────────────────────────────────────────────────────
  const exportAll = `${API_BASE}/export`;

  return (
    <AdminGuard>
      <div className="min-h-screen bg-bgLight relative overflow-hidden">
        <AnimatedBackground />

      <div className="max-w-7xl mx-auto px-4 py-10 relative z-10 space-y-6">

        {/* ── Header ─────────────────────────────────────────────────────── */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white/60 backdrop-blur-sm border border-white rounded-3xl p-6 shadow-sm">
          <div>
            <Breadcrumb items={[
              { label: "Home", href: "/" },
              { label: "Panel Admin", href: "/dashboard" },
              { label: "Interpretasi & Grafik" },
            ]} />
            <h1 className="text-2xl md:text-3xl font-extrabold text-textMain tracking-tight">
              Interpretasi &amp; Grafik Kuesioner
            </h1>
            <p className="text-slate-500 font-medium mt-1 text-sm">
              {data.length} partisipan terdaftar — SRQ-29 &amp; IPIP-BFM-50
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <a
              href={exportAll}
              download
              className="flex items-center gap-2 px-5 py-2.5 bg-palette5 text-white font-bold rounded-xl shadow-md hover:bg-palette4 transition-all text-sm"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Export Semua (CSV)
            </a>
          </div>
        </div>

        {/* ── Stat Cards ─────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-5">
          {[
            {
              label: "Rata-rata Neurotic",
              value: avgNeurotic.toFixed(1),
              note: "Indikasi GME jika ≥ 5",
              color: "text-palette4",
              borderColor: "border-palette4/20",
            },
            {
              label: "Total Partisipan",
              value: data.length.toString(),
              note: "Data terkumpul",
              color: "text-textMain",
              borderColor: "border-slate-100",
            },
            {
              label: "Alert Psikotik / PTSD",
              value: urgentCount.toString(),
              note: "Membutuhkan rujukan segera",
              color: "text-red-500",
              borderColor: "border-red-100",
            },
          ].map((stat, i) => (
            <motion.div
              key={stat.label}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.1 }}
              className={`bg-white rounded-2xl p-6 border ${stat.borderColor} shadow-sm`}
            >
              <p className="text-xs font-black tracking-widest text-slate-400 uppercase mb-2">{stat.label}</p>
              <p className={`text-5xl font-black ${stat.color}`}>{stat.value}</p>
              <p className="text-xs text-slate-400 font-medium mt-3 bg-slate-50 px-2 py-1 rounded-lg">{stat.note}</p>
            </motion.div>
          ))}
        </div>

        {/* ── Charts ─────────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* IPIP Radar */}
          <BentoCard>
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
              Rata-rata Profil IPIP-BFM-50
            </h2>
            <div className="h-64 select-none">
              <ResponsiveContainer width="100%" height="100%">
                <RadarChart cx="50%" cy="50%" outerRadius="75%" data={avgIPIP}>
                  <PolarGrid stroke="#e2e8f0" />
                  <PolarAngleAxis dataKey="subject" tick={{ fill: "#64748b", fontSize: 12, fontWeight: 700 }} />
                  <PolarRadiusAxis angle={30} domain={[0, 5]} tick={{ fill: "#94a3b8", fontSize: 10 }} />
                  <Radar name="Rata-rata" dataKey="value" stroke="#4a8fe7" fill="#59d2fe" fillOpacity={0.45} strokeWidth={2} />
                  <Tooltip
                    contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", fontSize: 13 }}
                    labelStyle={{ fontWeight: 700 }}
                  />
                </RadarChart>
              </ResponsiveContainer>
            </div>
            <p className="text-xs text-center text-slate-400 font-medium mt-3">Skala Mean 1–5 per dimensi</p>
          </BentoCard>

          {/* SRQ Bar */}
          <BentoCard>
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
              Skor Neurotic per Partisipan (SRQ-29)
            </h2>
            {srqBarData.length > 0 ? (
              <div className="h-64 select-none">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={srqBarData} margin={{ top: 5, right: 10, left: -20, bottom: 5 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" vertical={false} />
                    <XAxis dataKey="name" tick={{ fill: "#64748b", fontSize: 11 }} axisLine={false} tickLine={false} />
                    <YAxis domain={[0, 20]} tick={{ fill: "#94a3b8", fontSize: 11 }} axisLine={false} tickLine={false} />
                    <Tooltip
                      contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0" }}
                      cursor={{ fill: "#f8fafc" }}
                    />
                    <Bar dataKey="skor" name="Neurotic" radius={[6, 6, 0, 0]}>
                      {srqBarData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.skor >= 6 ? "#f87171" : entry.skor >= 5 ? "#fbbf24" : "#34d399"} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="h-64 flex items-center justify-center text-slate-400 font-medium text-sm">
                Belum ada data SRQ-29
              </div>
            )}
            <div className="flex gap-3 mt-3 text-xs font-bold">
              <span className="flex items-center gap-1.5"><span className="w-2.5 h-2.5 rounded-full bg-emerald-400 inline-block" />Normal (&lt;5)</span>
              <span className="flex items-center gap-1.5"><span className="w-2.5 h-2.5 rounded-full bg-amber-400 inline-block" />Indikasi (5)</span>
              <span className="flex items-center gap-1.5"><span className="w-2.5 h-2.5 rounded-full bg-red-400 inline-block" />GME (≥6)</span>
            </div>
          </BentoCard>
        </div>

        {/* ── Participant Table ───────────────────────────────────────────── */}
        <BentoCard>
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400">
              Detail Interpretasi Partisipan
            </h2>
            <span className="text-xs font-bold text-slate-400 bg-slate-50 px-3 py-1 rounded-full">{data.length} data</span>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-sm whitespace-nowrap">
              <thead className="bg-slate-50 text-slate-500 text-xs font-black uppercase tracking-wider">
                <tr>
                  <th className="px-4 py-3 text-left rounded-l-xl">Nama</th>
                  <th className="px-4 py-3 text-left">Neurotic</th>
                  <th className="px-4 py-3 text-left">Klasifikasi SRQ</th>
                  <th className="px-4 py-3 text-left">Dominan IPIP</th>
                  <th className="px-4 py-3 text-right rounded-r-xl">Ekspor</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {data.map((item, i) => {
                  const srq = item.score?.srq_score;
                  const severe = hasSevereFlag(item.score);
                  return (
                    <motion.tr
                      key={item.participant.id}
                      initial={{ opacity: 0, y: 4 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: i * 0.04 }}
                      className="hover:bg-slate-50 transition-colors"
                    >
                      <td className="px-4 py-3.5">
                        <div className="font-bold text-textMain">{item.participant.name}</div>
                        <div className="text-[10px] text-slate-400 font-mono">
                          {new Date(item.participant.created_at).toLocaleDateString("id-ID")}
                        </div>
                      </td>
                      <td className="px-4 py-3.5">
                        {srq ? (
                          <span className={`px-2.5 py-1 rounded-lg text-xs font-black ${srq.neurotic_score >= 6 ? "bg-red-100 text-red-600" : srq.neurotic_score >= 5 ? "bg-amber-100 text-amber-600" : "bg-emerald-100 text-emerald-700"}`}>
                            {srq.neurotic_score} / 20
                          </span>
                        ) : <span className="text-slate-300">—</span>}
                      </td>
                      <td className="px-4 py-3.5">
                        <span className={`text-xs font-bold ${severe ? "text-red-500" : (srq?.neurotic_score ?? 0) >= 5 ? "text-amber-500" : "text-emerald-600"}`}>
                          {getSRQClass(item.score)}
                        </span>
                      </td>
                      <td className="px-4 py-3.5">
                        <span className="text-xs font-bold text-palette5">{getDominantTrait(item.score)}</span>
                      </td>
                      <td className="px-4 py-3.5 text-right">
                        <a
                          href={`${API_BASE}/export/${item.participant.id}`}
                          download
                          className="text-xs font-bold text-slate-400 hover:text-palette5 underline decoration-dashed underline-offset-4 transition-colors"
                        >
                          CSV ↓
                        </a>
                      </td>
                    </motion.tr>
                  );
                })}
                {data.length === 0 && (
                  <tr>
                    <td colSpan={5} className="px-4 py-10 text-center text-slate-400 font-medium">
                      Belum ada partisipan yang menyelesaikan kuesioner.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </BentoCard>

      </div>
    </div>
    </AdminGuard>
  );
}
