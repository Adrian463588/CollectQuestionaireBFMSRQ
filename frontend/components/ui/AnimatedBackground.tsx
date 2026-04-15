

export function AnimatedBackground() {
  return (
    <>
      <div className="absolute top-[-20%] left-[-10%] w-[50%] h-[50%] bg-palette4/10 rounded-full blur-[120px] pointer-events-none" />
      <div className="absolute bottom-[-20%] right-[-10%] w-[50%] h-[50%] bg-palette3/10 rounded-full blur-[120px] pointer-events-none" />
    </>
  );
}
