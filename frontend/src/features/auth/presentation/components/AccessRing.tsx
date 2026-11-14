import { type RingSegment, SegmentedRing } from "../../../../components/ui/organisms/SegmentedRing";
import { t } from "../../../../lib/i18n/translate";

/** The shape of a month of charges, only to show what the app does (no user data yet). */
const sampleMonth: readonly RingSegment[] = [
  { value: 27, series: 1 },
  { value: 23, series: 2 },
  { value: 12, series: 3 },
  { value: 12, series: 4 },
  { value: 10, series: 5 },
  { value: 16, series: 6 },
];

/** The ring of the access screens (Registro, Login, Código). */
export function AccessRing() {
  return <SegmentedRing variant="normal" segments={sampleMonth} label={t("ui.access_ring")} size={120} />;
}
