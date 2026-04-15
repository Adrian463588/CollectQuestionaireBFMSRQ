"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar,
} from "recharts";
import { fetchAllDashboardData, DashboardItem, API_BASE } from "@/lib/api";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";
import { AdminGuard } from "@/components/AdminGuard";

// ── Helpers ───────────────────────────────────────────────────────────────────

function getSRQClass(s: DashboardItem["score"]): string {
  const srq = s?.srq_score;
  if (!srq) return "-";
  const flags: string[] = [];
  if (srq.neurotic_score >= 5) flags.push("Indikasi GME");
  if (srq.substance_use) flags.push("Penggunaan Zat");
  if (srq.psychotic) flags.push("Gejala Psikotik");
  if (srq.ptsd) flags.push("Gejala PTSD");
  return flags.length > 0 ? flags.join(" | ") : "Normal";
}

function getIPIPDominant(s: DashboardItem["score"]): string {
  const ipip = s?.ipip_score;
  if (!ipip) return "-";
  const dims = {
    Extraversion: ipip.extraversion,
    Agreeableness: ipip.agreeableness,
    Conscientiousness: ipip.conscientiousness,
    "Emotional Stability": ipip.emotional_stability,
    Intellect: ipip.intellect,
  };
  return Object.entries(dims).reduce((a, b) => (b[1] > a[1] ? b : a))[0];
}

// ── Component ─────────────────────────────────────────────────────────────────

export default function DashboardPage() {
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
      <div className="flex h-screen items-center justify-center bg-bgLight">
        <div className="animate-spin h-10 w-10 border-4 border-palette4 border-t-transparent rounded-full" />
      </div>
    );
  }

  // ── Derived stats ───────────────────────────────────────────────────────────
  const withSRQ = data.filter((d) => d.score?.srq_score);
  const withIPIP = data.filter((d) => d.score?.ipip_score);

  const avgNeurotic = withSRQ.length
    ? withSRQ.reduce((acc, d) => acc + (d.score!.srq_score!.neurotic_score ?? 0), 0) / withSRQ.length
    : 0;

  const urgentCount = data.filter((d) => d.score?.srq_score?.psychotic || d.score?.srq_score?.ptsd).length;

  const radarData = [
    { subject: "Extraversion", A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.extraversion ?? 0), 0) / withIPIP.length : 0, fullMark: 50 },
    { subject: "Agreeableness", A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.agreeableness ?? 0), 0) / withIPIP.length : 0, fullMark: 50 },
    { subject: "Conscientiousness", A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.conscientiousness ?? 0), 0) / withIPIP.length : 0, fullMark: 50 },
    { subject: "Emot. Stability", A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.emotional_stability ?? 0), 0) / withIPIP.length : 0, fullMark: 50 },
    { subject: "Intellect", A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.intellect ?? 0), 0) / withIPIP.length : 0, fullMark: 50 },
  ];

  const ipipBarData = withIPIP.map((d) => ({
    name: d.participant.name.split(" ")[0].slice(0, 8),
    E: d.score!.ipip_score!.extraversion,
    A: d.score!.ipip_score!.agreeableness,
    C: d.score!.ipip_score!.conscientiousness,
    S: d.score!.ipip_score!.emotional_stability,
    I: d.score!.ipip_score!.intellect,
  }));

  return (
    <AdminGuard>
      <div className="min-h-screen bg-bgLight relative overflow-hidden">
        <AnimatedBackground />

      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.35 }}
        className="max-w-7xl mx-auto px-4 py-10 relative z-10 space-y-6"
      >
        {/* ── Header ───────────────────────────────────────────────────────── */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white/60 backdrop-blur-sm border border-white rounded-3xl p-6 shadow-sm">
          <div>
            <h1 className="text-2xl md:text-3xl font-extrabold text-textMain tracking-tight">
              Admin Dashboard
            </h1>
            <p className="text-slate-500 font-medium mt-1 text-sm">
              Rekapitulasi {data.length} partisipan — SRQ-29 & IPIP-BFM-50
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <a
              href={`${API_BASE}/export`}
              download
              className="flex items-center gap-2 px-5 py-2.5 bg-palette5 text-white font-bold rounded-xl shadow-md hover:bg-palette4 transition-all text-sm"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Export Semua (CSV)
            </a>
            <Button variant="outline" onClick={() => window.location.href = "/interpretasi"}>
              Grafik Interpretasi
            </Button>
            <Button variant="secondary" onClick={() => window.location.href = "/"}>
              Home
            </Button>
          </div>
        </div>

        {/* ── Stat Cards ───────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-5">
          {[
            { label: "Rata-rata Neurotic", value: avgNeurotic.toFixed(1), note: "Indikasi GME jika ≥ 5", valueClass: "text-palette4" },
            { label: "Total Partisipan", value: String(data.length), note: "Data terkumpul", valueClass: "text-textMain" },
            { label: "Alert Urgent", value: String(urgentCount), note: "Psikotik / PTSD — rujukan segera", valueClass: "text-red-500" },
          ].map((stat, i) => (
            <motion.div
              key={stat.label}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.1 }}
              className="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm"
            >
              <p className="text-xs font-black tracking-widest uppercase text-slate-400 mb-2">{stat.label}</p>
              <p className={`text-5xl font-black ${stat.valueClass}`}>{stat.value}</p>
              <p className="text-xs text-slate-400 font-medium mt-3 bg-slate-50 px-2 py-1 rounded-lg">{stat.note}</p>
            </motion.div>
          ))}
        </div>

        {/* ── Charts ───────────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <BentoCard>
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
              Profil IPIP-BFM Rata-rata
            </h2>
            <div className="h-64 select-none">
              <ResponsiveContainer width="100%" height="100%">
                <RadarChart cx="50%" cy="50%" outerRadius="75%" data={radarData}>
                  <PolarGrid stroke="#e2e8f0" />
                  <PolarAngleAxis dataKey="subject" tick={{ fill: "#64748b", fontSize: 11, fontWeight: 700 }} />
                  <PolarRadiusAxis angle={30} domain={[0, 50]} tick={{ fill: "#94a3b8", fontSize: 10 }} />
                  <Radar name="Rata-rata" dataKey="A" stroke="#4a8fe7" fill="#59d2fe" fillOpacity={0.45} strokeWidth={2} />
                  <Tooltip contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", fontSize: 12 }} />
                </RadarChart>
              </ResponsiveContainer>
            </div>
          </BentoCard>

          <BentoCard>
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
              Distribusi Kepribadian per Partisipan
            </h2>
            {ipipBarData.length > 0 ? (
              <div className="h-64 select-none">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={ipipBarData} margin={{ top: 5, right: 10, left: -20, bottom: 5 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" vertical={false} />
                    <XAxis dataKey="name" tick={{ fill: "#64748b", fontSize: 11 }} axisLine={false} tickLine={false} />
                    <YAxis domain={[0, 50]} tick={{ fill: "#94a3b8", fontSize: 11 }} axisLine={false} tickLine={false} />
                    <Tooltip cursor={{ fill: "#f8fafc" }} contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", fontSize: 12 }} />
                    <Bar dataKey="E" fill="#73fbd3" stackId="a" radius={[0, 0, 4, 4]} name="Ext" />
                    <Bar dataKey="A" fill="#44e5e7" stackId="a" name="Agr" />
                    <Bar dataKey="C" fill="#59d2fe" stackId="a" name="Con" />
                    <Bar dataKey="S" fill="#4a8fe7" stackId="a" name="Sta" />
                    <Bar dataKey="I" fill="#5c7aff" stackId="a" radius={[4, 4, 0, 0]} name="Int" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="h-64 flex items-center justify-center text-slate-400 text-sm font-medium">
                Belum ada data IPIP-BFM-50
              </div>
            )}
          </BentoCard>
        </div>

        {/* ── Table ────────────────────────────────────────────────────────── */}
        <BentoCard>
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-sm font-black tracking-widest uppercase text-slate-400">
              Detail Skoring & Aksi
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
                  const severe = srq?.psychotic || srq?.ptsd || srq?.substance_use;
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
                        <span className={`text-xs font-bold ${severe ? "text-red-500" : srq?.neurotic_score! >= 5 ? "text-amber-500" : "text-emerald-600"}`}>
                          {getSRQClass(item.score)}
                        </span>
                      </td>
                      <td className="px-4 py-3.5">
                        <span className="text-xs font-bold text-palette5">{getIPIPDominant(item.score)}</span>
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
                      Belum ada partisipan yang terdaftar.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </BentoCard>
      </motion.div>
    </div>
    </AdminGuard>
  );
}
