# baggins

Utils to calculate sales values

## Install

```bash
$ go get github.com/profe-ajedrez/baggins
```

## Usage

Baggins is the handler provided to perform the sales operations over sales values

Internally Baggins has a handler for taxes and a handler for discounts which
performs operations and calculations over these concepts.

Baggins can register different types of taxes and discounts and is able to
calculate them correctly.


### Taxes

Baggins uses the concept of stages to the taxes registry and calculations,
where  a tax can be registered in a particular stage which determines when is calculated.

#### The taxes stages are:

  * OverTaxableStage   represents taxes calculated over its value.

  * OverTaxesStage represents taxes calculated over its value plus the cummulated amount of the taxes calculated in the OvertaxableStage

  * OverTaxesIgnorableStage represents taxes which are calculated like the taxes of the OverTaxableStage, but are not included in the OVerTaxesStage

```go
b := baggins.New()

// adds a percentual tax to the Overtaxable stage
err  := b.AddTax(decimal.NewFromInt(10), tax.PercentualMode, tax.OverTaxableStage)

if err != nil {
    panic(err) // Remember! Dont Panic!
}
```

### Discounts 

You can register discounts in baggins.

```go

b := baggins.New()

// register a percentual discount
err := b.AddDiscount(decimal.NewFromInt(10), discount.PercentualMode)

if err != nil {
    panic(err) // Remember! Dont Panic!
}
```


### Calculate results

When you are done registering taxes an discount you can invoke the method `Calculate`.

```go
b := baggins.New()

// adds a percentual tax to the Overtaxable stage
err  := b.AddTax(decimal.NewFromInt(10), tax.PercentualMode, tax.OverTaxableStage)

if err != nil {
    panic(err) // Remember! Dont Panic!
}

// register a percentual discount
err := b.AddDiscount(decimal.NewFromInt(10), discount.PercentualMode)

if err != nil {
    panic(err) // Remember! Dont Panic!
}

unitValue := decimal.NewFromInt(100)
quantity := decimal.NewFromInt(10)

result, err := b.Calculate(unitValue, quantity)

if err != nil {
    panic(err) // Remember! Dont Panic!
}

js, _ := json.Marshal(result)
fmt.Println(js)

// prints
//
//  {{"withDiscount":{"net":"900","brute":"990","tax":"90","discount":"10","discountedValue":"100","discountedValueBrute":"110","unitValue":"90"},"withoutDiscount":{"net":"1000","brute":"1100","tax":"100","unitValue":"100"}}

```