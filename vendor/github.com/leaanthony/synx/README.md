# Synx

Simple wrappers to allow concurrent writes to scalar types.

## Install

```go get github.com/leaanthony/synx```

## Usage

Each supported type has a New<Type> function, EG NewInt(). This takes an intial value. 
Accessing the value is simply through GetValue() and SetValue().

### New<type>

Create a new Synx wrapper for type. 

```
 var a = synx.NewInt(0)
```

### GetValue()

Return the current value of the scalar.

```
 var b := a.GetValue()
```

## SetValue()

Update the current value of the scalar.

```
 var c := a.SetValue(100)
```

## Lock() / Unlock()

Enables / Disables the synchronisation locks for the scalar

```
 c.Lock()
 // Perform operations requiring exclusive lock 
 c.Unlock()
```

