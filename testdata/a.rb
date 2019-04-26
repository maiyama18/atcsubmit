a = 5.times.map { gets.to_i }
k = gets.to_i

def can_communicate(a, k)
  (0..4).each do |i|
    (0..4).each do |j|
      return false if (a[i] - a[j]).abs > k
    end
  end
  true
end

puts can_communicate(a, k) ? 'Yay!' : ':('
