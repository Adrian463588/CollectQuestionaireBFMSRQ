"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { saveParticipantToSession } from "@/lib/participant";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";

export default function Home() {
  const router = useRouter();
  const [participantName, setParticipantName] = useState("");
  const [age, setAge] = useState("");
  const [gender, setGender] = useState<"male" | "female">("male");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!participantName || !age) {
      alert("Please fill in all fields");
      return;
    }

    saveParticipantToSession({
      name: participantName,
      age: parseInt(age),
      gender,
      participantId: undefined,
    });

    router.push("/questionnaire-select");
  };

  return (
    <main className="min-h-screen bg-bgLight relative flex items-center justify-center p-4 overflow-hidden">
      <AnimatedBackground />

      <BentoCard className="max-w-md w-full">
        <div className="mb-8 text-center">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: "spring", stiffness: 200, delay: 0.2 }}
            className="w-16 h-16 bg-gradient-to-tr from-palette4 to-palette3 rounded-2xl mx-auto mb-4 flex items-center justify-center shadow-lg shadow-palette4/30"
          >
            <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg>
          </motion.div>
          <h1 className="text-3xl font-extrabold text-textMain mb-2 tracking-tight">
            Kuesioner Psikologi
          </h1>
          <p className="text-textMuted font-medium">IPIP-BFM-50 & SRQ-29</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="space-y-1.5">
            <label className="block text-textMain text-sm font-bold pl-1">
              Nama Lengkap
            </label>
            <input
              type="text"
              value={participantName}
              onChange={(e) => setParticipantName(e.target.value)}
              className="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 text-textMain placeholder-slate-400 font-medium focus:outline-none focus:border-palette4 focus:ring-4 focus:ring-palette4/15 transition-all shadow-sm"
              placeholder="Masukkan nama Anda"
              required
            />
          </div>

          <div className="space-y-1.5">
            <label className="block text-textMain text-sm font-bold pl-1">
              Usia
            </label>
            <input
              type="number"
              value={age}
              onChange={(e) => setAge(e.target.value)}
              className="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 text-textMain placeholder-slate-400 font-medium focus:outline-none focus:border-palette4 focus:ring-4 focus:ring-palette4/15 transition-all shadow-sm"
              placeholder="Contoh: 20"
              min="15"
              required
            />
          </div>

          <div className="space-y-1.5">
            <label className="block text-textMain text-sm font-bold pl-1 mb-2">
              Jenis Kelamin
            </label>
            <div className="flex gap-3">
              <Button 
                type="button"
                variant="likert"
                isActive={gender === "male"} 
                onClick={() => setGender("male")}
                className="flex-1"
              >
                Laki-laki
              </Button>
              <Button 
                type="button" 
                variant="likert"
                isActive={gender === "female"} 
                onClick={() => setGender("female")}
                className="flex-1"
              >
                Perempuan
              </Button>
            </div>
          </div>

          <Button type="submit" variant="primary" className="w-full mt-6 group">
            Mulai Kuesioner
            <svg className="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
          </Button>
        </form>

        <div className="mt-8 text-center pt-6 border-t border-slate-100 flex items-center justify-center gap-5">
          <a
            href="/interpretasi"
            className="text-slate-400 hover:text-palette4 text-sm transition-colors font-bold flex items-center gap-1 group"
          >
            Lihat Grafik Interpretasi
            <svg className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5" d="M9 5l7 7-7 7" />
            </svg>
          </a>
          <span className="text-slate-200">|</span>
          <a
            href="/dashboard"
            className="text-slate-400 hover:text-slate-600 text-sm transition-colors font-semibold"
          >
            Panel Admin
          </a>
        </div>
      </BentoCard>
    </main>
  );
}
