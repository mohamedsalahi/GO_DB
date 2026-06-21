import { motion } from 'framer-motion';
import { ListTodo, CheckCircle2, Clock, TrendingUp } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: typeof ListTodo;
  delay?: number;
}

const icons = {
  ListTodo, CheckCircle2, Clock, TrendingUp,
};

export function StatsCard({ title, value, icon: Icon, delay = 0 }: StatsCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.25, delay }}
      className="bg-surface-card border border-zinc-800 rounded-xl p-4"
    >
      <div className="flex items-center justify-between mb-1.5">
        <p className="text-2xs text-zinc-600 font-medium uppercase tracking-wider">{title}</p>
        <Icon className="w-3.5 h-3.5 text-zinc-600" strokeWidth={1.5} />
      </div>
      <p className="text-xl font-semibold text-zinc-100 tabular-nums">{value}</p>
    </motion.div>
  );
}
