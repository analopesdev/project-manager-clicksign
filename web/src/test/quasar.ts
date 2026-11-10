import {
  ClosePopup,
  QBtn,
  QCard,
  QDate,
  QIcon,
  QImg,
  QInput,
  QItem,
  QList,
  QMenu,
  QPopupProxy,
  Quasar,
  Ripple
} from "quasar";

export const quasarPlugin: [typeof Quasar, Record<string, unknown>] = [
  Quasar,
  {
    components: {
      QBtn,
      QCard,
      QDate,
      QIcon,
      QImg,
      QInput,
      QItem,
      QList,
      QMenu,
      QPopupProxy
    },
    directives: { ClosePopup, Ripple }
  }
];
