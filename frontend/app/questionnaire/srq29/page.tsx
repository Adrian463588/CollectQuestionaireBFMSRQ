"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { AnimatePresence, motion } from "framer-motion";
import { srq29Questions } from "@/data/srq29";
import { SRQAnswer } from "@/types";
import { createParticipant, submitQuestionnaire } from "@/lib/api";
import { getParticipantFromSession, saveParticipantIdToSession } from "@/lib/participant";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { ProgressBar } from "@/components/ui/ProgressBar";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";
import { Breadcrumb } from "@/components/ui/Breadcrumb";

export default function SRQ29Questionnaire() {
  const router = useRouter();
  const [currentQuestion, setCurrentQuestion] = useState(0);
  const [answers, setAnswers] = useState<SRQAnswer>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleAnswer = (answer: boolean) => {
    const questionId = srq29Questions[currentQuestion].id.toString();
    setAnswers((prev) => ({ ...prev, [questionId]: answer }));
    if (currentQuestion < srq29Questions.length - 1) {
      setCurrentQuestion((q) => q + 1);
    }
  };

  const handlePrevious = () => {
    if (currentQuestion > 0) setCurrentQuestion((q) => q - 1);
  };

  const handleSubmit = async () => {
    setIsSubmitting(true);
    try {
      const participant = getParticipantFromSession();
      if (!participant) { router.push("/"); return; }

      let participantId = participant.participantId;
      if (!participantId) {
        const created = await createParticipant({ name: participant.name, age: participant.age, gender: participant.gender });
        participantId = created.id;
        saveParticipantIdToSession(participantId);
      }

      await submitQuestionnaire(participantId, "srq29", answers);
      router.push(`/results/srq29?participantId=${participantId}`);
    } catch (error) {
      console.error("Error submitting:", error);
      alert("Gagal mengirim jawaban, coba lagi.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const currentQ = srq29Questions[currentQuestion];
  const answeredCount = Object.keys(answers).length;
  const progress = (answeredCount / srq29Questions.length) * 100;
  const isFinished = answeredCount === srq29Questions.length;
  const currentAnswered = answers[currentQ?.id?.toString()];

  return (
    <main className="min-h-screen bg-bgLight relative flex items-center justify-center p-4 overflow-hidden">
      <AnimatedBackground />

      <BentoCard className="max-w-2xl w-full">
        <Breadcrumb items={[
          { label: "Home", href: "/" },
          { label: "Pilih Kuesioner", href: "/questionnaire-select" },
          { label: "SRQ-29" },
        ]} />
        {/* Header */}
        <div className="mb-8">
          <div className="flex justify-between items-center mb-4">
            <span className="text-palette4 font-black tracking-widest text-xs uppercase bg-palette4/10 px-3 py-1.5 rounded-full">
              SRQ-29
            </span>
            <span className="text-slate-400 font-semibold text-sm tabular-nums">
              {currentQuestion + 1} / {srq29Questions.length}
            </span>
          </div>
          <ProgressBar progress={progress} />
        </div>

        {/* Question */}
        <div className="min-h-[140px] flex items-center justify-center mb-8 relative">
          <AnimatePresence mode="wait">
            <motion.div
              key={currentQuestion}
              initial={{ opacity: 0, x: 30 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -30 }}
              transition={{ duration: 0.25, ease: "easeOut" }}
              className="absolute w-full"
            >
              <h3 className="text-xl md:text-2xl font-bold text-textMain text-center leading-relaxed">
                {currentQ?.text}
              </h3>
            </motion.div>
          </AnimatePresence>
        </div>

        {/* Answer buttons */}
        <div className="grid grid-cols-2 gap-4 mb-8">
          <Button
            variant="likert"
            isActive={currentAnswered === true}
            onClick={() => handleAnswer(true)}
            className="flex-col gap-1.5 py-7"
          >
            <span className="text-2xl">✓</span>
            <span className="font-bold text-base">Ya</span>
          </Button>
          <Button
            variant="likert"
            isActive={currentAnswered === false}
            onClick={() => handleAnswer(false)}
            className="flex-col gap-1.5 py-7"
          >
            <span className="text-2xl">✗</span>
            <span className="font-bold text-base">Tidak</span>
          </Button>
        </div>

        {/* Navigation */}
        <div className="flex justify-between items-center pt-5 border-t border-slate-100">
          <Button
            variant="outline"
            onClick={handlePrevious}
            disabled={currentQuestion === 0 || isSubmitting}
            className="text-sm px-5 py-2.5 gap-1.5"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 19l-7-7 7-7" />
            </svg>
            Kembali
          </Button>

          {isFinished && currentQuestion === srq29Questions.length - 1 ? (
            <Button variant="primary" onClick={handleSubmit} disabled={isSubmitting}>
              {isSubmitting ? (
                <span className="flex items-center gap-2">
                  <svg className="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  Memproses...
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  Lihat Hasil
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                </span>
              )}
            </Button>
          ) : (
            <Button
              variant="secondary"
              onClick={() => setCurrentQuestion((q) => q + 1)}
              disabled={currentAnswered === undefined || currentQuestion >= srq29Questions.length - 1}
              className="px-5 py-2.5 text-sm gap-1.5"
            >
              Lanjut
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
              </svg>
            </Button>
          )}
        </div>
      </BentoCard>
    </main>
  );
}
