const defaultTx = {
  UUID: "00000000-0000-0000-0000-000000000000",
  DateCreated: "0001-01-01T00:00:00",
  DateModified: "0001-01-01T00:00:00",
  Date: undefined,
  Description: "",
  Amount: "",
  Shared: false,
  GeoLocation: "",
  User: {
    ID: 1
  },
  Method: {
    ID: 1
  },
  Category: {
    ID: 0
  },
  Type: {
    ID: 0
  }
};

var vm = new Vue({
  el: '#app',
  data: {
    editMode: false,
    transaction: {},
    users: [],
    categories: [],
    types: [],
    transactions: [],
    methods: []
  },
  computed: {
    SharedQuota: function () {
      if (Array.isArray(this.transaction.Shares)) {
        return this.transaction.Shares.reduce((sum, shr) => {
          const q = new Decimal(shr.Quota);

          return q.plus(sum)
        }, new Decimal(0));
      } else return "";
    },
    OwnQuota: function () {
      if (Array.isArray(this.transaction.Shares)) {
        const a = new Decimal(this.transaction.Amount);
        const sq = new Decimal(this.SharedQuota);
        return a.minus(sq).toString();
      } else return "";
    },
    TxDate: {
      get: function () {
        return this.justDate(this.transaction.Date);
      },
      set: function (v) {
        console.log("yo");
        this.transaction.Date = v + "T00:00:00";
      }
    }


  },
  methods: {
    fetchInitialState: function () {
      fetch('/api/home/')
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => {
          this.transactions = data.transactions;
          this.types = data.types;
          this.users = data.users;
          this.categories = data.categories;
          //this.transaction = data.transaction;
          this.methods = data.methods;
        })
        .catch(res => console.log(res))

    },
    fetchLatest: function (n=5, offset=0,orderBy="date_modified DESC, date DESC") {
      fetch('/api/transactions/?limit=' + n + '&offset=' + offset + '&orderBy=' + orderBy)
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => this.transactions.push(...data))
        .catch(res => console.log(res))
    },
    reloadLatest: function () {
      n = this.transactions.length;
      this.transactions = [];
      this.fetchLatest(n);
    },
    more: function (e) {
      e.preventDefault();
      this.fetchLatest(5, this.transactions.length);
    },
    justDate: function (dateStr = "") {

      if (dateStr == "") {
        const date = new Date();
        dateStr = date.toISOString();
      }
      return dateStr.split('T')[0];
    },
    fetchTransaction: function (uuid, e) {
      e.preventDefault();
      fetch('/api/transaction/' + uuid)
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => {
          this.editMode = true;
          this.transaction = data;
        })
        .catch(res => console.log(res))

    },
    deleteTx: function (e) {
      e.preventDefault();

      const uuid = this.transaction.UUID;

      fetch('/api/transaction/' + uuid, {
          method: `DELETE`,
          credentials: `same-origin`
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          const i = this.transactions.findIndex(function (t) {
            if (t.UUID == uuid) return true
          });

          this.transactions.splice(i, 1)
        })
        .catch(res => console.log(res))
    },
    updateTx: function (e) {
      e.preventDefault();
      fetch('/api/transaction/', {
          method: `PUT`,
          credentials: `same-origin`,
          body: JSON.stringify(this.transaction)
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          this.reloadLatest();
          this.editMode = false;
          this.transaction = defaultTx;//We have to copy the object!!!!!!!!!
        })
        .catch(res => console.log(res))

    },
    addTx: function (e) {
      e.preventDefault();
      fetch('/api/transaction/', {
          method: `POST`,
          credentials: `same-origin`,
          body: JSON.stringify(this.transaction)
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          this.reloadLatest();
          this.transaction = defaultTx;
        })
        .catch(res => console.log(res))

    },
    addShare: function (e) {
      e.preventDefault();
      if (!this.transaction.Shared) {
        return;
      }

      if (!Array.isArray(this.transaction.Shares)) {
        this.transaction.Shares = [];
      }

      this.transaction.Shares.push({
        WithID: 0,
        Quota: 0
      });
    },

  },
  created: function () {
    this.fetchInitialState();
    this.transaction = defaultTx;
  }
})