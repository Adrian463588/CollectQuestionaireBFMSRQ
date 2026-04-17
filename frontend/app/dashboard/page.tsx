"use client";

import { useEffect, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, RadarChart, PolarGrid, PolarAngleAxis,
  PolarRadiusAxis, Radar,
} from "recharts";
import { fetchAllDashboardData, deleteParticipant, DashboardItem, API_BASE } from "@/lib/api";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";
import { AdminGuard } from "@/components/AdminGuard";
import { Breadcrumb } from "@/components/ui/Breadcrumb";

// ── Label helpers ─────────────────────────────────────────────────────────────

/** Human-readable label for SRQ overall_risk */
function riskLabel(risk: string): { text: string; cls: string } {
  switch (risk) {
    case "kritis":  return { text: "Kritis",  cls: "bg-red-100 text-red-700" };
    case "tinggi":  return { text: "Tinggi",  cls: "bg-orange-100 text-orange-700" };
    case "sedang":  return { text: "Sedang",  cls: "bg-amber-100 text-amber-700" };
    default:        return { text: "Rendah",  cls: "bg-emerald-100 text-emerald-700" };
  }
}

/** Human-readable label for IPIP dimension interpretation */
function ipipLabel(label: string): string {
  const map: Record<string, string> = {
    sangat_tinggi: "Sangat Tinggi",
    tinggi:        "Tinggi",
    rata_rata:     "Rata-rata",
    rendah:        "Rendah",
    sangat_rendah: "Sangat Rendah",
  };
  return map[label] ?? label;
}

/** Returns the dominant IPIP trait name based on mean score */
function getIPIPDominant(item: DashboardItem["score"]): string {
  const ipip = item?.ipip_score;
  if (!ipip) return "—";
  const dims: Record<string, number> = {
    Extraversion:        ipip.extraversion,
    Agreeableness:       ipip.agreeableness,
    Conscientiousness:   ipip.conscientiousness,
    "Emot. Stability":   ipip.emotional_stability,
    Intellect:           ipip.intellect,
  };
  return Object.entries(dims).reduce((a, b) => (b[1] > a[1] ? b : a))[0];
}

// ── Delete Confirmation Modal ─────────────────────────────────────────────────

interface DeleteModalProps {
  participantName: string;
  onConfirm: () => Promise<void>;
  onCancel: () => void;
}

function DeleteModal({ participantName, onConfirm, onCancel }: DeleteModalProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState<string | null>(null);

  const handleConfirm = async () => {
    setLoading(true);
    setError(null);
    try {
      await onConfirm();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Gagal menghapus partisipan");
      setLoading(false);
    }
  };

  return (
    <motion.div
      key="overlay"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4"
      onClick={onCancel}
    >
      <motion.div
        key="modal"
        initial={{ opacity: 0, scale: 0.9, y: 16 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.9, y: 16 }}
        transition={{ type: "spring", stiffness: 380, damping: 28 }}
        className="bg-white rounded-2xl shadow-2xl p-8 max-w-md w-full"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Icon */}
        <div className="flex items-center justify-center w-14 h-14 bg-red-100 rounded-2xl mx-auto mb-5">
          <svg className="w-7 h-7 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>

        <h2 className="text-xl font-extrabold text-gray-900 text-center mb-2">
          Hapus Partisipan?
        </h2>
        <p className="text-sm text-gray-500 text-center mb-1">
          Data berikut akan dihapus dari sistem:
        </p>
        <p className="text-base font-bold text-gray-800 text-center mb-6 bg-gray-50 rounded-xl py-2 px-4">
          {participantName}
        </p>

        <p className="text-xs text-gray-400 text-center mb-6 leading-relaxed">
          Seluruh jawaban dan hasil skoring akan ikut terhapus.
          Data tetap tersimpan di database untuk keperluan audit riset.
        </p>

        {error && (
          <p className="text-xs text-red-500 text-center mb-4 font-medium">{error}</p>
        )}

        <div className="flex gap-3">
          <button
            onClick={onCancel}
            disabled={loading}
            className="flex-1 py-2.5 rounded-xl border border-gray-200 text-gray-600 font-bold text-sm
                       hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onClick={handleConfirm}
            disabled={loading}
            className="flex-1 py-2.5 rounded-xl bg-red-600 text-white font-bold text-sm
                       hover:bg-red-700 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {loading ? (
              <span className="h-4 w-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            ) : null}
            {loading ? "Menghapus…" : "Ya, Hapus"}
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
}

// ── Main Component ────────────────────────────────────────────────────────────

export default function DashboardPage() {
  const [data, setData]             = useState<DashboardItem[]>([]);
  const [loading, setLoading]       = useState(true);
  const [deleteTarget, setDeleteTarget] = useState<DashboardItem | null>(null);

  const loadData = useCallback(() => {
    setLoading(true);
    fetchAllDashboardData()
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadData();
  }, [loadData]);

  // Optimistic deletion: remove from local state, then call API.
  const handleDelete = useCallback(async () => {
    if (!deleteTarget) return;
    const target = deleteTarget;
    // Optimistic update
    setData((prev) => prev.filter((d) => d.participant.id !== target.participant.id));
    setDeleteTarget(null);
    await deleteParticipant(target.participant.id);
    // Full refresh to sync aggregate stats
    loadData();
  }, [deleteTarget, loadData]);

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-bgLight">
        <div className="animate-spin h-10 w-10 border-4 border-palette4 border-t-transparent rounded-full" />
      </div>
    );
  }

  // ── Derived stats ─────────────────────────────────────────────────────────
  const withSRQ  = data.filter((d) => d.score?.srq_score);
  const withIPIP = data.filter((d) => d.score?.ipip_score);

  const avgNeurotic = withSRQ.length
    ? withSRQ.reduce((acc, d) => acc + (d.score!.srq_score!.neurotic_score ?? 0), 0) / withSRQ.length
    : 0;

  const urgentCount = data.filter(
    (d) => d.score?.srq_score?.psychotic || d.score?.srq_score?.ptsd,
  ).length;

  // ── Chart data (mean scale 1–5) ───────────────────────────────────────────
  const radarData = [
    { subject: "Extraversion",       A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.extraversion ?? 0), 0)         / withIPIP.length : 0, fullMark: 5 },
    { subject: "Agreeableness",      A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.agreeableness ?? 0), 0)        / withIPIP.length : 0, fullMark: 5 },
    { subject: "Conscientiousness",  A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.conscientiousness ?? 0), 0)    / withIPIP.length : 0, fullMark: 5 },
    { subject: "Emot. Stability",    A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.emotional_stability ?? 0), 0)  / withIPIP.length : 0, fullMark: 5 },
    { subject: "Intellect",          A: withIPIP.length ? withIPIP.reduce((a, d) => a + (d.score!.ipip_score!.intellect ?? 0), 0)            / withIPIP.length : 0, fullMark: 5 },
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

        {/* ── Delete confirmation modal ──────────────────────────────────── */}
        <AnimatePresence>
          {deleteTarget && (
            <DeleteModal
              participantName={deleteTarget.participant.name}
              onConfirm={handleDelete}
              onCancel={() => setDeleteTarget(null)}
            />
          )}
        </AnimatePresence>

        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35 }}
          className="max-w-7xl mx-auto px-4 py-10 relative z-10 space-y-6"
        >
          <Breadcrumb items={[
            { label: "Home", href: "/" },
            { label: "Panel Admin" },
          ]} />
          {/* ── Header ─────────────────────────────────────────────────── */}
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4
                          bg-white/60 backdrop-blur-sm border border-white rounded-3xl p-6 shadow-sm">
            <div>
              <h1 className="text-2xl md:text-3xl font-extrabold text-textMain tracking-tight">
                Admin Dashboard
              </h1>
              <p className="text-slate-500 font-medium mt-1 text-sm">
                Rekapitulasi {data.length} partisipan — SRQ-29 &amp; IPIP-BFM-50
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <a
                href={`${API_BASE}/export`}
                download
                className="flex items-center gap-2 px-5 py-2.5 bg-palette5 text-white font-bold
                           rounded-xl shadow-md hover:bg-palette4 transition-all text-sm"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                    d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Export Semua (CSV)
              </a>
              <Button variant="outline" onClick={() => (window.location.href = "/interpretasi")}>
                Grafik Interpretasi
              </Button>
              <Button variant="secondary" onClick={() => (window.location.href = "/")}>
                Home
              </Button>
            </div>
          </div>

          {/* ── Stat Cards ─────────────────────────────────────────────── */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-5">
            {[
              { label: "Rata-rata Neurotic", value: avgNeurotic.toFixed(1),   note: "Indikasi GME jika ≥ 5",            cls: "text-palette4" },
              { label: "Total Partisipan",  value: String(data.length),       note: "Data terkumpul",                   cls: "text-textMain" },
              { label: "Alert Urgent",       value: String(urgentCount),       note: "Psikotik / PTSD — rujukan segera", cls: "text-red-500"  },
            ].map((stat, i) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm"
              >
                <p className="text-xs font-black tracking-widest uppercase text-slate-400 mb-2">{stat.label}</p>
                <p className={`text-5xl font-black ${stat.cls}`}>{stat.value}</p>
                <p className="text-xs text-slate-400 font-medium mt-3 bg-slate-50 px-2 py-1 rounded-lg">{stat.note}</p>
              </motion.div>
            ))}
          </div>

          {/* ── Charts ─────────────────────────────────────────────────── */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <BentoCard>
              <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
                Profil IPIP-BFM Rata-rata (Skala 1–5)
              </h2>
              <div className="h-64 select-none">
                <ResponsiveContainer width="100%" height="100%">
                  <RadarChart cx="50%" cy="50%" outerRadius="75%" data={radarData}>
                    <PolarGrid stroke="#e2e8f0" />
                    <PolarAngleAxis dataKey="subject" tick={{ fill: "#64748b", fontSize: 11, fontWeight: 700 }} />
                    <PolarRadiusAxis angle={30} domain={[0, 5]} tick={{ fill: "#94a3b8", fontSize: 10 }} />
                    <Radar name="Rata-rata" dataKey="A" stroke="#4a8fe7" fill="#59d2fe" fillOpacity={0.45} strokeWidth={2} />
                    <Tooltip
                      contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", fontSize: 12 }}
                      formatter={(v) => [typeof v === "number" ? v.toFixed(2) : v, "Mean"]}
                    />
                  </RadarChart>
                </ResponsiveContainer>
              </div>
            </BentoCard>

            <BentoCard>
              <h2 className="text-sm font-black tracking-widest uppercase text-slate-400 mb-5">
                Distribusi Kepribadian per Partisipan (Mean 1–5)
              </h2>
              {ipipBarData.length > 0 ? (
                <div className="h-64 select-none">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={ipipBarData} margin={{ top: 5, right: 10, left: -20, bottom: 5 }}>
                      <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" vertical={false} />
                      <XAxis dataKey="name" tick={{ fill: "#64748b", fontSize: 11 }} axisLine={false} tickLine={false} />
                      <YAxis domain={[0, 5]} tick={{ fill: "#94a3b8", fontSize: 11 }} axisLine={false} tickLine={false} />
                      <Tooltip
                        cursor={{ fill: "#f8fafc" }}
                        contentStyle={{ borderRadius: 12, border: "1px solid #e2e8f0", fontSize: 12 }}
                        formatter={(v) => [typeof v === "number" ? v.toFixed(2) : v]}
                      />
                      <Bar dataKey="E" fill="#73fbd3" stackId="a" radius={[0, 0, 4, 4]} name="Extraversion" />
                      <Bar dataKey="A" fill="#44e5e7" stackId="a" name="Agreeableness" />
                      <Bar dataKey="C" fill="#59d2fe" stackId="a" name="Conscientiousness" />
                      <Bar dataKey="S" fill="#4a8fe7" stackId="a" name="Emot. Stability" />
                      <Bar dataKey="I" fill="#5c7aff" stackId="a" radius={[4, 4, 0, 0]} name="Intellect" />
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

          {/* ── Table ──────────────────────────────────────────────────── */}
          <BentoCard>
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-sm font-black tracking-widest uppercase text-slate-400">
                Detail Skoring &amp; Aksi
              </h2>
              <span className="text-xs font-bold text-slate-400 bg-slate-50 px-3 py-1 rounded-full">
                {data.length} data
              </span>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-sm whitespace-nowrap">
                <thead className="bg-slate-50 text-slate-500 text-xs font-black uppercase tracking-wider">
                  <tr>
                    <th className="px-4 py-3 text-left rounded-l-xl">Nama</th>
                    <th className="px-4 py-3 text-left">Neurotic</th>
                    <th className="px-4 py-3 text-left">Risiko Keseluruhan</th>
                    <th className="px-4 py-3 text-left">Dominan IPIP</th>
                    <th className="px-4 py-3 text-left">Interpretasi IPIP</th>
                    <th className="px-4 py-3 text-center">Ekspor</th>
                    <th className="px-4 py-3 text-center rounded-r-xl">Hapus</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  <AnimatePresence mode="popLayout">
                    {data.map((item, i) => {
                      const srq       = item.score?.srq_score;
                      const ipip      = item.score?.ipip_score;
                      const risk      = srq ? riskLabel(srq.overall_risk) : null;
                      const dominant  = getIPIPDominant(item.score);
                      const domLabel  = ipip
                        ? (() => {
                            const map: Record<string, string> = {
                              Extraversion:       ipip.extra_label,
                              Agreeableness:      ipip.agre_label,
                              Conscientiousness:  ipip.cons_label,
                              "Emot. Stability":  ipip.stab_label,
                              Intellect:          ipip.intell_label,
                            };
                            return ipipLabel(map[dominant] ?? "");
                          })()
                        : null;

                      return (
                        <motion.tr
                          key={item.participant.id}
                          layout
                          initial={{ opacity: 0, y: 4 }}
                          animate={{ opacity: 1, y: 0 }}
                          exit={{ opacity: 0, x: -20 }}
                          transition={{ delay: i * 0.03 }}
                          className="hover:bg-slate-50 transition-colors"
                        >
                          {/* Name */}
                          <td className="px-4 py-3.5">
                            <div className="font-bold text-textMain">{item.participant.name}</div>
                            <div className="text-[10px] text-slate-400 font-mono">
                              {new Date(item.participant.created_at).toLocaleDateString("id-ID")}
                            </div>
                          </td>

                          {/* Neurotic score */}
                          <td className="px-4 py-3.5">
                            {srq ? (
                              <span className={`px-2.5 py-1 rounded-lg text-xs font-black ${
                                srq.neurotic_score >= 6 ? "bg-red-100 text-red-600"
                                : srq.neurotic_score >= 5 ? "bg-amber-100 text-amber-600"
                                : "bg-emerald-100 text-emerald-700"
                              }`}>
                                {srq.neurotic_score} / 20
                              </span>
                            ) : <span className="text-slate-300">—</span>}
                          </td>

                          {/* Overall Risk */}
                          <td className="px-4 py-3.5">
                            {risk ? (
                              <span className={`text-xs font-black px-2.5 py-1 rounded-lg ${risk.cls}`}>
                                {risk.text}
                              </span>
                            ) : <span className="text-slate-300">—</span>}
                          </td>

                          {/* Dominant IPIP trait */}
                          <td className="px-4 py-3.5">
                            <span className="text-xs font-bold text-palette5">{dominant}</span>
                          </td>

                          {/* IPIP label for dominant trait */}
                          <td className="px-4 py-3.5">
                            {domLabel
                              ? <span className="text-xs text-slate-600 font-medium">{domLabel}</span>
                              : <span className="text-slate-300">—</span>}
                          </td>

                          {/* Export CSV */}
                          <td className="px-4 py-3.5 text-center">
                            <a
                              href={`${API_BASE}/export/${item.participant.id}`}
                              download
                              className="text-xs font-bold text-slate-400 hover:text-palette5
                                         underline decoration-dashed underline-offset-4 transition-colors"
                            >
                              CSV ↓
                            </a>
                          </td>

                          {/* Delete button */}
                          <td className="px-4 py-3.5 text-center">
                            <button
                              id={`delete-${item.participant.id}`}
                              onClick={() => setDeleteTarget(item)}
                              className="inline-flex items-center justify-center w-8 h-8 rounded-lg
                                         text-slate-400 hover:text-red-500 hover:bg-red-50
                                         transition-colors"
                              aria-label={`Hapus ${item.participant.name}`}
                            >
                              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </td>
                        </motion.tr>
                      );
                    })}
                  </AnimatePresence>

                  {data.length === 0 && (
                    <tr>
                      <td colSpan={7} className="px-4 py-10 text-center text-slate-400 font-medium">
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
