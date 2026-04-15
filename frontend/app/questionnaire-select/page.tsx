'use client';

import { useRouter } from "next/navigation";
import { getParticipantFromSession } from "@/lib/participant";
import { useEffect, useState } from "react";

import { BentoCard } from "@/components/ui/BentoCard";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";

export default function QuestionnaireSelect() {
  const router = useRouter();
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    const participant = getParticipantFromSession();
    if (!participant) {
      router.push("/");
    } else {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsLoaded(true);
    }
  }, [router]);

  const handleSelect = (type: "srq29" | "ipip") => {
    router.push(`/questionnaire/${type}`);
  };

  if (!isLoaded) return null;

  return (
    <main className="min-h-screen bg-bgLight relative flex items-center justify-center p-4 overflow-hidden">
      <AnimatedBackground />

      <div className="max-w-4xl w-full relative z-10">
        <div className="mb-12 text-center">
          <span className="bg-palette5/10 text-palette5 font-bold tracking-widest text-xs px-3 py-1 rounded-full uppercase">Kustomisasi Pilihan</span>
          <h1 className="text-4xl md:text-5xl font-extrabold text-textMain mt-4 tracking-tight">
            Pilih Kuesioner
          </h1>
          <p className="text-textMuted font-medium mt-3">Selesaikan satu persatu untuk menyelesaikan asesmen Anda.</p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          <BentoCard hoverEffect onClick={() => handleSelect('srq29')} delay={0.1}>
            <div className="w-14 h-14 bg-palette4/10 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-palette4 transition-all">
              <svg className="w-7 h-7 text-palette4 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            </div>
            <h2 className="text-3xl font-extrabold text-textMain mb-3 group-hover:text-palette4 transition-colors tracking-tight">SRQ-29</h2>
            <p className="text-textMuted font-medium mb-6 leading-relaxed">
              Self Reporting Questionnaire untuk skrining kesehatan mental dasar.
            </p>
            <ul className="space-y-3">
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">29 pertanyaan singkat</span>
              </li>
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">Respons Ya/Tidak</span>
              </li>
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">Skrining Psikotik & Narkotik</span>
              </li>
            </ul>
          </BentoCard>

          <BentoCard hoverEffect onClick={() => handleSelect('ipip')} delay={0.2}>
            <div className="w-14 h-14 bg-palette5/10 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-palette5 transition-all">
              <svg className="w-7 h-7 text-palette5 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            </div>
            <h2 className="text-3xl font-extrabold text-textMain mb-3 group-hover:text-palette5 transition-colors tracking-tight">IPIP-BFM-50</h2>
            <p className="text-textMuted font-medium mb-6 leading-relaxed">
              International Personality Item Pool - Pengukuran kepribadian menyeluruh.
            </p>
            <ul className="space-y-3">
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">50 pernyataan profil</span>
              </li>
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">1-5 Skala Likert Setuju/Tidak</span>
              </li>
              <li className="flex items-start gap-3">
                <div className="w-5 h-5 mt-0.5 rounded-full bg-palette1/20 flex items-center justify-center flex-shrink-0 text-palette1"><svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"></path></svg></div>
                <span className="text-textMain font-medium text-sm">Pemetaan 5 Dimensi Kepribadian</span>
              </li>
            </ul>
          </BentoCard>
        </div>
      </div>
    </main>
  );
}
